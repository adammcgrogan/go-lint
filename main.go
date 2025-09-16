package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("Error: No file specified. Usage: go run . <filename.go>")
	}
	fileName := os.Args[1]

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Running linter on %s...\n", fileName)

	ast.Inspect(node, func(n ast.Node) bool {
		if cfg.Rules.CheckExportedComments {
			checkExportedComments(fset, n)
		}
		if cfg.Rules.CheckMagicStrings.Enabled {
			checkMagicStrings(fset, n, cfg.Rules.CheckMagicStrings.MaxLength)
		}
		if cfg.Rules.CheckParameterCount.Enabled {
			checkParameterCount(fset, n, cfg.Rules.CheckParameterCount.Max)
		}
		if cfg.Rules.CheckFunctionLength.Enabled {
			checkFunctionLength(fset, n, cfg.Rules.CheckFunctionLength.MaxLines)
		}
		return true
	})

	fmt.Println("Linting complete.")

}
