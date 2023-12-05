package main

import (
	"fmt"
	"go/ast"

	"github.com/jingyugao/rowserrcheck/passes/rowserr"
	"github.com/timakin/bodyclose/passes/bodyclose"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// OsExitCheckAnalyzer checks for os.Exit call in main function of main package
var OsExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osexitcheck",
	Doc:  "check for os.Exit call in main function of main package",
	Run:  run,
}

var staticcheckAnalyzers = map[string]bool{
	// Various misuses of the standard library
	"SA1*": true,
	// Concurrency issues
	"SA2*": true,
	// Testing issues
	"SA3*": true,
	// Code that isn't really doing anything
	"SA4*": true,
	// Correctness issues
	"SA5*": true,
	// Performance issues
	"SA6*": true,
	// 	Dubious code constructs that have a high probability of being wrong
	"SA9*": true,
}

var stylecheckAnalyzers = map[string]bool{
	// Incorrectly formatted error string
	"ST1005": true,
	// Poorly chosen receiver name
	"ST1006": true,
}

var customAnalyzers = []*analysis.Analyzer{
	// check for os.Exit call in main function of main package
	OsExitCheckAnalyzer,
}

var passesAnalyzers = []*analysis.Analyzer{
	// Checks consistency of Printf format strings and arguments
	printf.Analyzer,
	// Checks for shadowed variables
	shadow.Analyzer,
	// Checks struct field tags are well formed
	structtag.Analyzer,
}

var publicAnalyzers = []*analysis.Analyzer{
	// checks whether HTTP response body is closed successfully
	bodyclose.Analyzer,
	// checks whether Rows.Err is checked
	rowserr.NewAnalyzer(
		"github.com/jmoiron/sqlx",
	),
}

type checks []*analysis.Analyzer

// adds staticcheck analyzers
func (c *checks) addStaticcheckAnalyzers() {
	for _, v := range staticcheck.Analyzers {
		// wildecard name check
		if staticcheckAnalyzers[fmt.Sprintf("%s*", v.Analyzer.Name[:3])] {
			*c = append(*c, v.Analyzer)
		}

		// specific name check
		if staticcheckAnalyzers[v.Analyzer.Name] {
			*c = append(*c, v.Analyzer)
		}
	}
}

// adds stylecheck analyzers
func (c *checks) addStylecheckAnalyzers() {
	for _, v := range stylecheck.Analyzers {
		// wildecard name check
		if stylecheckAnalyzers[fmt.Sprintf("%s*", v.Analyzer.Name[:3])] {
			*c = append(*c, v.Analyzer)
		}

		// specific name check
		if stylecheckAnalyzers[v.Analyzer.Name] {
			*c = append(*c, v.Analyzer)
		}
	}
}

// add passes analyzers
func (c *checks) addPassesAnalyzers() {
	*c = append(*c, passesAnalyzers...)
}

// add custom analyzers
func (c *checks) addCustomAnalyzers() {
	*c = append(*c, customAnalyzers...)
}

// add custom analyzers
func (c *checks) addPublicAnalyzers() {
	*c = append(*c, publicAnalyzers...)
}

func run(pass *analysis.Pass) (any, error) {
	for _, f := range pass.Files {
		// обходим дерево разбора
		ast.Inspect(f, func(n ast.Node) bool {
			// интересует только пакет main
			if p, ok := n.(*ast.File); ok && p.Name.Name == "main" {

				// смотрим какие функции есть в пакете
				for _, decl := range p.Decls {
					// интересует только функция main
					if mainFunc, ok := decl.(*ast.FuncDecl); ok && mainFunc.Name.Name == "main" {

						// смотрим что там в ней внутри
						for _, b := range mainFunc.Body.List {
							// интересуют только вызовы функций
							if exp, ok := b.(*ast.ExprStmt); ok {
								if caller, ok := exp.X.(*ast.CallExpr); ok {
									// ищем функию Exit
									if exitFunc, ok := caller.Fun.(*ast.SelectorExpr); ok {
										// именно из пакета os
										if pkg, ok := exitFunc.X.(*ast.Ident); ok && pkg.Name == "os" {
											if exitFunc.Sel.Name == "Exit" {
												pass.Reportf(exitFunc.Sel.NamePos, "calling os.Exit in main function is not allowed")
											}
										}
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}

func addAnalyzers(c *checks) {
	c.addCustomAnalyzers()
	c.addStaticcheckAnalyzers()
	c.addStylecheckAnalyzers()
	c.addPassesAnalyzers()
	c.addPublicAnalyzers()
}

func main() {
	var c checks
	addAnalyzers(&c)

	multichecker.Main(
		c...,
	)
}
