package regexponce_test

import (
	"testing"

	"github.com/budougumi0617/regexponce"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, regexponce.Analyzer, "a")
}

func Test_targetFuncs(t *testing.T) {
	got, err := regexponce.TargetFuncs()
	if err != nil {
		t.Errorf("targetFuncs() error = %v", err)
		return
	}
	if len(got) != 4 {
		t.Errorf("targetFuncs() returns %v", got)
	}
}
