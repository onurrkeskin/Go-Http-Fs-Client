package analyzer

type SingleCharAnalyzer struct {
	selectedChar byte
}

func NewSingleCharAnalyzer(verifierOpts ...string) *SingleCharAnalyzer {
	return &SingleCharAnalyzer{selectedChar: verifierOpts[0][0]}
}

func (runeVerifier *SingleCharAnalyzer) Analyze(analyzerData AnalyzerData) AnalyzerResult {
	result := AnalyzerResult{
		Position: -1,
	}

	for in, curRune := range analyzerData.CurData {
		if curRune == runeVerifier.selectedChar {
			result.Position = in
			return result
		}
	}
	return result
}
