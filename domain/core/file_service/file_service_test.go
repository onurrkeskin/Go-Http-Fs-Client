package file_service_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"gitlab.com/onurkeskin/go-http-fs-client/domain/core/file_service"
	"gitlab.com/onurkeskin/go-http-fs-client/domain/util/testinghelpers"
	"gitlab.com/onurkeskin/go-http-fs-client/foundation/docker"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

var c *docker.DockerContainer

func TestMain(m *testing.M) {
	var err error
	c, err = testinghelpers.StartFS()
	if err != nil {
		fmt.Println(err)
		return
	}

	m.Run()
}

func TestFilesService(t *testing.T) {
	fileService := file_service.NewFileService(c.Host, "./")

	t.Log("Given the current server with multiple z word containing files")
	{
		testID := 0
		retrieveFilesResponse, err := fileService.GetFiles(context.Background(), "z")
		if err != nil {
			t.Fatalf("\t%s\tTest %d:\tShouldnt have thrown errors: %s.", Failed, testID, err)
		}
		if strings.Contains(retrieveFilesResponse.DownloadedFiles, "file4") && strings.Contains(retrieveFilesResponse.DownloadedFiles, "file2") {
			t.Logf("\t%s\tTest %d:\tFound them on 2nd and 4th files", Success, testID)
		} else {
			t.Fatalf("\t%s\tTest %d:\tShould have found those on 2nd and 4th files : %s.", Failed, testID, err)

		}
	}
}
