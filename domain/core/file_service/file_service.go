package file_service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"gitlab.com/onurkeskin/go-http-fs-client/domain/core/file_service/analyzer"
	"go.uber.org/zap"
)

const (
	FAILED_DL = ""
)

var (
	ErrFilesNotFound = errors.New("No Matching files found")
	ErrMustBeARune   = errors.New("Ther are more than single characters to look for")
)

var (
	gofileServerFileIndicatorRegexp = regexp.MustCompile(`<a href=".+">(.+)<\/a>`)
)

type FileService struct {
	logger              *zap.SugaredLogger
	fileServerUrl       string
	DownloadLocation    string
	curAnalyzerStrategy analyzer.ANALYZER_STRATEGY
}

func NewFileService(fileServerUrl string, downloadLocation string) FileService {
	return FileService{
		fileServerUrl:       fileServerUrl,
		DownloadLocation:    downloadLocation,
		curAnalyzerStrategy: analyzer.ANALYZER_STRATEGY_CHAR,
	}
}

func (fileService FileService) GetFiles(ctx context.Context, significantCharacter string) (RetrieveFilesResponse, error) {
	resp, err := http.Get("http://" + fileService.fileServerUrl)
	if err != nil {
		return RetrieveFilesResponse{}, errors.New("Server unreachable" + fileService.fileServerUrl)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return RetrieveFilesResponse{}, errors.New("Cant read http response" + fileService.fileServerUrl)
	}
	foundRegexGroups := gofileServerFileIndicatorRegexp.FindAllStringSubmatch(string(body), -1)
	foundFiles := make([]string, len(foundRegexGroups))
	for in, group := range foundRegexGroups {
		if len(group) > 0 {
			foundFiles[in] = group[1]
		}
	}

	curAnalyzer, err := fileService.newAnalyzer(significantCharacter)
	if err != nil {
		return RetrieveFilesResponse{}, errors.New("Unkown analyzer configuration" + fileService.fileServerUrl)
	}

	var counterLock sync.RWMutex
	var sharedCounter int = math.MaxInt
	var barrier sync.WaitGroup
	waitChannel := make(chan struct{})
	errorChannel := make(chan error, len(foundFiles))
	successfulDownloadsTrackerChannel := make(chan string, len(foundFiles))
	completedDownloadsChannel := make(chan string)
	// defer close(completedDownloadsChannel)
	barrier.Add(len(foundFiles))
	for _, file := range foundFiles {
		// probably some sort of pooling would be nice according to maxprocs
		go fileService.downloadFile(ctx, FileDownloaderOpts{
			url:      fileService.fileServerUrl,
			fileName: file,
			maxSize:  1000,
			az:       curAnalyzer,
		},
			FileDownloaderConcurrencyHelpers{
				sharedPositionLock:  &counterLock,
				sharedFoundPosition: &sharedCounter,
				barrier:             &barrier,
			},
			FileDownloaderListeners{
				errorChannel:      errorChannel,
				successfulChannel: successfulDownloadsTrackerChannel,
			},
		)
	}

	go func() {
		defer close(waitChannel)
		barrier.Wait()
	}()

	go func() {
		defer close(successfulDownloadsTrackerChannel)
		completedFiles := strings.Builder{}
		for i := 0; i < len(foundFiles); i++ {
			if file := <-successfulDownloadsTrackerChannel; file != "" {
				completedFiles.WriteString(fmt.Sprintf("Completed: %s ", file))
			}
		}
		completedDownloadsChannel <- completedFiles.String()
	}()

FileDownloadProcess:
	for {
		select {
		case e := <-errorChannel:
			fileService.logger.Error(e)
		case <-waitChannel:
			close(errorChannel)
			break FileDownloadProcess
		}
	}

	return RetrieveFilesResponse{DownloadedFiles: <-completedDownloadsChannel}, nil
}

func (fileService FileService) newAnalyzer(verifierOpts ...string) (analyzer.Analyzer, error) {
	switch fileService.curAnalyzerStrategy {
	case analyzer.ANALYZER_STRATEGY_CHAR:
		return analyzer.NewSingleCharAnalyzer(verifierOpts...), nil
	default:
		return nil, errors.New("Invalid Analyzer Configuration")
	}
}

type FileDownloaderOpts struct {
	url      string
	fileName string
	maxSize  int64
	az       analyzer.Analyzer
}

type FileDownloaderConcurrencyHelpers struct {
	sharedPositionLock  *sync.RWMutex
	sharedFoundPosition *int
	barrier             *sync.WaitGroup
}

type FileDownloaderListeners struct {
	errorChannel      chan error
	successfulChannel chan string
}

func (fileService *FileService) downloadFile(ctx context.Context, data FileDownloaderOpts, fdConcurrentOpts FileDownloaderConcurrencyHelpers, fdListener FileDownloaderListeners) {
	fileLocation := fileService.DownloadLocation + "/" + data.fileName
	conn, err := net.Dial("tcp", data.url)
	if err != nil {
		fdListener.errorChannel <- fmt.Errorf("Cant connect to file server to retrieve file: %s", fileLocation)
		fdListener.successfulChannel <- FAILED_DL
		fdConcurrentOpts.barrier.Done()
		return
	}
	defer conn.Close()

	fo, err := os.Create(fileLocation)
	if err != nil {
		fdListener.errorChannel <- fmt.Errorf("Cant create file for: %s", fileLocation)
		fdConcurrentOpts.barrier.Done()
		fdListener.successfulChannel <- FAILED_DL
		return
	}
	defer fo.Close()
	conn.Write([]byte(fmt.Sprintf("GET /%s HTTP/1.0\r\nHost: localhost\r\n\r\n", data.fileName)))

	buffer := make([]byte, data.maxSize)
	curPos := 0
	selfFoundPosition := math.MaxInt
	for n, err := conn.Read(buffer); n > 0 && err == nil; n, err = conn.Read(buffer) {
		result := data.az.Analyze(analyzer.AnalyzerData{
			FHandle: fo,
			CurData: buffer,
		})

		// If a match has been found, stop and wait until others finished
		if result.Position > -1 {
			selfFoundPosition = curPos + result.Position
			fdConcurrentOpts.sharedPositionLock.Lock()
			if *fdConcurrentOpts.sharedFoundPosition >= selfFoundPosition {
				fo.Write(buffer)
				*fdConcurrentOpts.sharedFoundPosition = selfFoundPosition
			}
			fdConcurrentOpts.sharedPositionLock.Unlock()
			break
		}

		// If a match has not been found, but there exists a file that has a match and it is on a prior byte than this routines current immediatly try to exit this routine
		curPos += n
		fdConcurrentOpts.sharedPositionLock.RLock()
		if curPos > *fdConcurrentOpts.sharedFoundPosition {
			if selfFoundPosition > *fdConcurrentOpts.sharedFoundPosition {
				os.Remove(fileLocation)
				fdConcurrentOpts.sharedPositionLock.RUnlock()
				fdConcurrentOpts.barrier.Done()
				fdListener.successfulChannel <- FAILED_DL
				return
			}
			// no need to continue further
			fdConcurrentOpts.sharedPositionLock.RUnlock()
			break
		}
		fdConcurrentOpts.sharedPositionLock.RUnlock()

		_, err := fo.Write(buffer)
		if err != nil {
			fdListener.errorChannel <- fmt.Errorf("Cant write to file: %s", fileLocation)
			fdConcurrentOpts.barrier.Done()
			fdListener.successfulChannel <- FAILED_DL
			return
		}
	}

	fdConcurrentOpts.barrier.Done()
	fdConcurrentOpts.barrier.Wait()

	if selfFoundPosition > *fdConcurrentOpts.sharedFoundPosition || selfFoundPosition == math.MaxInt {
		fdListener.successfulChannel <- FAILED_DL
		os.Remove(fileLocation)
	} else {
		// download the rest of the file since its safe
		buffer := make([]byte, 1024)
		for n, err := conn.Read(buffer); err == nil && n > 0; n, err = conn.Read(buffer) {
			fo.Write(buffer[0:n])
		}
		fdListener.successfulChannel <- fileLocation
	}
}
