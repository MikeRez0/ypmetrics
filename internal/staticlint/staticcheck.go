package staticlint

import (
	"golang.org/x/tools/go/analysis"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func GetStaticCheck() []*analysis.Analyzer {
	res := make([]*analysis.Analyzer, 0, 10)
	for _, sa := range staticcheck.Analyzers {
		res = append(res, sa.Analyzer)
	}

	return res
}

func GetStyleCheck() []*analysis.Analyzer {
	res := make([]*analysis.Analyzer, 0, 10)
	for _, sa := range stylecheck.Analyzers {
		res = append(res, sa.Analyzer)
	}

	return res
}
