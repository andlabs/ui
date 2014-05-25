// 24 may 2014
package main

import (
	"fmt"
	"os"
	"strings"
	"go/token"
	"go/ast"
	"go/parser"
)

func main() {
	var pkg *ast.Package

	importpath := os.Args[1]

	fileset := token.NewFileSet()

	filter := func(i os.FileInfo) bool {
		return strings.HasSuffix(i.Name(), "_windows.go")
	}
	pkgs, err := parser.ParseDir(fileset, importpath,
		filter, parser.AllErrors)
	if err != nil {
		panic(err)
	}
	if len(pkgs) != 1 {
		panic("more than one package found")
	}
	for k, _ := range pkgs {		// get the sole key
		pkg = pkgs[k]
	}

	var run func(...ast.Decl)
	var runstmt func(ast.Stmt)
	var runblock func(*ast.BlockStmt)

	desired := func(name string) bool {
		return strings.HasPrefix(name, "_")
	}
	run = func(decls ...ast.Decl) {
		for _, d := range decls {
			switch dd := d.(type) {
			case *ast.FuncDecl:
				runblock(dd.Body)
			case *ast.GenDecl:
				if desired(d.Name.String()) {
					fmt.Println(d.Name.String())
				}
			default:
				panic(fmt.Errorf("unknown decl type %T: %v", dd, dd))
			}
		}
	}
	runstmt = func(s ast.Stmt) {
		switch ss := s.(type) {
		case *ast.DeclStmt:
			run(ss.Decl)
		case *ast.LabeledStmt:
			runstmt(ss.Stmt)
		case *ast.AssignStmt:
			// TODO go through Lhs if ss.Tok type == DEFINE
		case *ast.GoStmt:
		// these don't have decls
		case *ast.EmptyStmt:
		case *ast.ExprStmt:
		case *ast.SendStmt:
		case *ast.IncDecStmt:
			// all do nothing
		default:
			panic(fmt.Errorf("unknown stmt type %T: %v", dd, dd))
		}
	}
	runblock = func(block *ast.BlockStmt) {
		for _, s := range block.Stmt {
			runstmt(s)
		}
	}
	for _, f := range pkg.Files {
		run(f.Decls...)
	}
}
