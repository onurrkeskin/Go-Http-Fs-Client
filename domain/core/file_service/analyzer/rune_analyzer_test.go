package analyzer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestRuneAnalyzerFinalChar(t *testing.T) {
	singleCharAnalyzer := NewSingleCharAnalyzer("a")

	testID := 0
	t.Log("Given a string of arbitrary characters with rune a on 8th position")
	{
		result := singleCharAnalyzer.Analyze(AnalyzerData{
			FHandle: nil,
			CurData: []byte(string("iopiopqa")),
		})
		expected := result
		expected.Position = 7
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tShould have found it on 7th position. Diff:\n%s", Failed, testID, diff)
		}
		t.Logf("\t%s\tTest %d:\tFound it on 8th position", Success, testID)
	}
}
func TestRuneAnalyzerFirstChar(t *testing.T) {
	singleCharAnalyzer := NewSingleCharAnalyzer("a")

	testID := 0
	t.Log("Given a string of arbitrary characters with rune a on 1st position")
	{
		result := singleCharAnalyzer.Analyze(AnalyzerData{
			FHandle: nil,
			CurData: []byte(string("aiopiopqa")),
		})
		expected := result
		expected.Position = 0
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tShould have found it on 1st position. Diff:\n%s", Failed, testID, diff)
		}

		t.Logf("\t%s\tTest %d:\tFound it on 8st position", Success, testID)
	}
}

func TestRuneAnalyzerMultipleArbitraryChar(t *testing.T) {
	singleCharAnalyzer := NewSingleCharAnalyzer("a")

	testID := 0
	t.Log("Given a string of arbitrary characters with rune a on 3rd position")
	{
		result := singleCharAnalyzer.Analyze(AnalyzerData{
			FHandle: nil,
			CurData: []byte(string("ioapiaopqa")),
		})
		expected := result
		expected.Position = 2
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tShould have found it on 3rd position. Diff:\n%s", Failed, testID, diff)
		}
		t.Logf("\t%s\tTest %d:\tFound it on 3rd position", Success, testID)
	}
}

func TestRuneAnalyzerNotFoundChar(t *testing.T) {
	singleCharAnalyzer := NewSingleCharAnalyzer("a")

	testID := 0
	t.Log("Given a string of arbitrary characters with no a in any position")
	{
		result := singleCharAnalyzer.Analyze(AnalyzerData{
			FHandle: nil,
			CurData: []byte(string("qweiopnm")),
		})
		expected := result
		expected.Position = -1
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Fatalf("\t%s\tTest %d:\tFound it on a position. Diff:\n%s", Failed, testID, diff)
		}
		t.Logf("\t%s\tTest %d:\t Not Found", Success, testID)
	}
}
