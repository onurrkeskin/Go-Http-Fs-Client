package analyzer

import "os"

type ANALYZER_STRATEGY int32

const (
	ANALYZER_STRATEGY_CHAR = ANALYZER_STRATEGY(1)
)

type AnalyzerData struct {
	FHandle *os.File
	CurData []byte
}

type AnalyzerResult struct {
	Position int
}

type Analyzer interface {
	Analyze(analyzerData AnalyzerData) AnalyzerResult
}

type AnalyzerFunc func(AnalyzerData) AnalyzerResult

func (analyzerFunc AnalyzerFunc) Analyze(analyzerData AnalyzerData) AnalyzerResult {
	return analyzerFunc(analyzerData)
}
