package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

func checkExportedComments(fset *token.FileSet, node ast.Node) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	if fn.Name.IsExported() && fn.Doc == nil {
		pos := fset.Position(fn.Pos())
		fmt.Printf("-> [comment-check] Exported function '%s' at %s is missing a comment.\n", fn.Name.Name, pos)
	}
}

func checkMagicStrings(fset *token.FileSet, node ast.Node, maxLength int) {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return
	}

	if lit.Kind == token.STRING {
		stringValue := lit.Value[1 : len(lit.Value)-1]
		if len(stringValue) > maxLength { // Use the configurable max length
			pos := fset.Position(lit.Pos())
			fmt.Printf("-> [magic-string] Found a long hardcoded string %s at %s. Consider using a constant.\n", lit.Value, pos)
		}
	}
}
