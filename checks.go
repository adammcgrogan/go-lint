package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
)

func collectNoLintLines(fset *token.FileSet, node *ast.File) map[int]bool {
	noLintLines := make(map[int]bool)
	for _, commentGroup := range node.Comments {
		for _, comment := range commentGroup.List {
			if strings.Contains(comment.Text, "nolint") {
				line := fset.Position(comment.Pos()).Line
				noLintLines[line] = true
			}
		}
	}
	return noLintLines
}

func checkExportedComments(fset *token.FileSet, node ast.Node, noLintLines map[int]bool) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	if fn.Name.IsExported() && fn.Doc == nil {
		line := fset.Position(fn.Pos()).Line
		if !noLintLines[line] {
			pos := fset.Position(fn.Pos())
			fmt.Printf("-> [comment-check] Exported function '%s' at %s is missing a comment.\n", fn.Name.Name, pos)
		}
	}
}

func checkMagicStrings(fset *token.FileSet, node ast.Node, maxLength int, noLintLines map[int]bool) {
	lit, ok := node.(*ast.BasicLit)
	if !ok {
		return
	}

	if lit.Kind == token.STRING {
		stringValue := lit.Value[1 : len(lit.Value)-1]
		if len(stringValue) > maxLength {
			line := fset.Position(lit.Pos()).Line
			if !noLintLines[line] {
				pos := fset.Position(lit.Pos())
				fmt.Printf("-> [magic-string] Found a long hardcoded string %s at %s. Consider using a constant.\n", lit.Value, pos)
			}
		}
	}
}

func checkParameterCount(fset *token.FileSet, node ast.Node, maxParams int, noLintLines map[int]bool) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	numParams := 0
	for _, field := range fn.Type.Params.List {
		numParams += len(field.Names)
	}

	if numParams > maxParams {
		line := fset.Position(fn.Pos()).Line
		if !noLintLines[line] {
			pos := fset.Position(fn.Pos())
			fmt.Printf("-> [param-count] Function '%s' at %s has %d parameters, which exceeds the max of %d.\n", fn.Name.Name, pos, numParams, maxParams)
		}
	}
}

func checkFunctionLength(fset *token.FileSet, node ast.Node, maxLines int, noLintLines map[int]bool) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	startPos := fset.Position(fn.Body.Pos())
	endPos := fset.Position(fn.Body.End())

	lineCount := endPos.Line - startPos.Line
	if lineCount > maxLines {
		declarationLine := fset.Position(fn.Pos()).Line

		if !noLintLines[declarationLine] {
			fmt.Printf("-> [func-length] Function '%s' at %s has %d lines, which exceeds the max of %d.\n", fn.Name.Name, startPos, lineCount, maxLines)
		}
	}
}

func checkDeferInLoop(fset *token.FileSet, node ast.Node, noLintLines map[int]bool) {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return
	}

	ast.Inspect(forStmt.Body, func(n ast.Node) bool {
		if deferStmt, ok := n.(*ast.DeferStmt); ok {
			line := fset.Position(deferStmt.Pos()).Line
			if !noLintLines[line] {
				pos := fset.Position(deferStmt.Pos())
				fmt.Printf("-> [defer-in-loop] Found a 'defer' statement inside a loop at %s. It will not execute until the function returns.\n", pos)
			}
		}
		return true
	})
}

type ReceiverNameChecker struct {
	fset          *token.FileSet
	receiverNames map[string]map[string]token.Pos
}

func NewReceiverNameChecker(fset *token.FileSet) *ReceiverNameChecker {
	return &ReceiverNameChecker{
		fset:          fset,
		receiverNames: make(map[string]map[string]token.Pos),
	}
}

func (c *ReceiverNameChecker) Visit(node ast.Node) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
		return
	}

	receiver := fn.Recv.List[0]
	receiverName := receiver.Names[0].Name

	var typeName string
	if starExpr, ok := receiver.Type.(*ast.StarExpr); ok {
		typeName = starExpr.X.(*ast.Ident).Name
	} else if ident, ok := receiver.Type.(*ast.Ident); ok {
		typeName = ident.Name
	}

	if typeName != "" {
		if _, ok := c.receiverNames[typeName]; !ok {
			c.receiverNames[typeName] = make(map[string]token.Pos)
		}
		c.receiverNames[typeName][receiverName] = receiver.Pos()
	}
}

func (c *ReceiverNameChecker) Report() {
	for structName, names := range c.receiverNames {

		if len(names) > 1 {
			fmt.Printf("-> [receiver-name] Inconsistent receiver names for struct '%s':\n", structName)
			for name, pos := range names {
				fmt.Printf("    - Found '%s' at %s\n", name, c.fset.Position(pos))
			}
		}
	}
}
