// staticlint app allows running checks on the source code.
// The following checks will run:
// `shadow.Analyzer` - check for possible unintended shadowing of variables
// `assign.Analyzer` - check for useless assignments
// `analyzers.ExitAnalyzer` - forbids using os.Exit() in the main package
// `simple.Analyzers` - Use plain channel send or receive instead of single-case select
package main

import (
	"github.com/naneri/shortener/cmd/staticlint/analyzers"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	var mychecks []*analysis.Analyzer

	mychecks = append(mychecks, shadow.Analyzer)
	mychecks = append(mychecks, assign.Analyzer)
	mychecks = append(mychecks, analyzers.ExitAnalyzer)
	mychecks = append(mychecks, simple.Analyzers["S1000"])

	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
		mychecks = append(mychecks, v)
	}

	multichecker.Main(
		mychecks...,
	)
}
