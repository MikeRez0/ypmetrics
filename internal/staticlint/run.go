// Package staticlint - custom code checker.
package staticlint

import (
	fatcontext "github.com/Crocmagnon/fatcontext/pkg/analyzer"
	"github.com/breml/errchkjson"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

// Run - functions runs analyzers.
//
// Included:
//
// + Standard passes
//
// + staticcheck - SA class
//
// + staticcheck - style checks
//
// + fatcontext Go linter which detects potential fat contexts in loops or function literals.
// They can lead to performance issues, as documented here: https://gabnotes.org/fat-contexts/
//
// + errchkjson Checks types passed to the json encoding functions.
// Reports unsupported types and reports occurrences where the check for the returned error can be omitted.
func Run() {
	// Создайте свой multichecker, состоящий из:
	mychecks := []*analysis.Analyzer{}

	// стандартных статических анализаторов пакета golang.org/x/tools/go/analysis/passes;
	mychecks = append(mychecks, GetStdPasses()...)

	// всех анализаторов класса SA пакета staticcheck.io;
	mychecks = append(mychecks, GetStaticCheck()...)

	// не менее одного анализатора остальных классов пакета staticcheck.io;
	mychecks = append(mychecks, GetStyleCheck()...)

	// двух или более любых публичных анализаторов на ваш выбор.
	mychecks = append(mychecks,
		fatcontext.Analyzer,
		errchkjson.NewAnalyzer(),
		// собственный анализатор, запрещающий использовать прямой вызов os.Exit в функции main пакета main
		OsExitCheck,
	)

	multichecker.Main(mychecks...)
}
