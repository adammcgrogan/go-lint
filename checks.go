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

func checkParameterCount(fset *token.FileSet, node ast.Node, maxParams int) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	numParams := len(fn.Type.Params.List)
	if numParams > maxParams {
		pos := fset.Position(fn.Pos())
		fmt.Printf("-> [param-count] Function '%s' at %s has %d parameters, which exceeds the max of %d.\n", fn.Name.Name, pos, numParams, maxParams)
	}
}

func checkFunctionLength(fset *token.FileSet, node ast.Node, maxLines int) {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return
	}

	startPos := fset.Position(fn.Body.Pos())
	endPos := fset.Position(fn.Body.End())

	lineCount := endPos.Line - startPos.Line
	if lineCount > maxLines {
		fmt.Printf("-> [func-length] Function '%s' at %s has %d lines, which exceeds the max of %d.\n", fn.Name.Name, startPos, lineCount, maxLines)
	}
}

func checkDeferInLoop(fset *token.FileSet, node ast.Node) {
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		return
	}

	ast.Inspect(forStmt.Body, func(n ast.Node) bool {
		if deferStmt, ok := n.(*ast.DeferStmt); ok {
			pos := fset.Position(deferStmt.Pos())
			fmt.Printf("-> [defer-in-loop] Found a 'defer' statement inside a loop at %s. It will not execute until the function returns.\n", pos)
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
