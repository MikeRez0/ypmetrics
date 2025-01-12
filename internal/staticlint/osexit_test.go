package staticlint_test

import (
	"testing"

	"github.com/MikeRez0/ypmetrics/internal/staticlint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOsExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), staticlint.OsExitCheck, "./...")
}