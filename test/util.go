package test

import (
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func MustParse(t *testing.T, fileSet *token.FileSet, file string, data []byte) *ast.File {
	f, err := parser.ParseFile(fileSet, file, data, parser.ParseComments)
	FailOnError(t, err)

	return f
}

func MustReadFile(t *testing.T, file string) []byte {
	b, err := ioutil.ReadFile(file)
	FailOnError(t, err)

	return b
}
func FailOnError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func GetFunctionBody(t *testing.T, data []byte, name string) string {
	var (
		fileSet = token.NewFileSet()
		f       = MustParse(t, fileSet, "", data)
	)
	sb := &strings.Builder{}

	ast.Inspect(f, func(node ast.Node) bool {
		if fd, ok := node.(*ast.FuncDecl); ok {
			if fd.Name.Name == name {
				FailOnError(t, printer.Fprint(sb, fileSet, fd.Body.List))
			}
		}
		return true
	})

	return sb.String()
}

func GetImports(t *testing.T, data []byte) string {
	var (
		fileSet = token.NewFileSet()
		f       = MustParse(t, fileSet, "", data)
	)
	sb := &strings.Builder{}

	ast.Inspect(f, func(node ast.Node) bool {
		if fd, ok := node.(*ast.File); ok {
			if d, ok1 := fd.Decls[0].(*ast.GenDecl); ok1 {
				if d.Tok != token.IMPORT {
					return true
				}
				FailOnError(t, printer.Fprint(sb, fileSet, d))
			}
		}
		return true
	})

	return sb.String()
}

func AssertEquals(t *testing.T, want string, got string) {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(want, got, false)

	if len(diffs) > 1 || (len(diffs) == 1 && diffs[0].Type != diffmatchpatch.DiffEqual) {
		t.Fatal(dmp.DiffPrettyText(diffs))
	}
}
