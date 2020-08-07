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

