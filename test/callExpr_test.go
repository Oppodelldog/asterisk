package test

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"testing"

	. "github.com/Oppodelldog/asterisk"
)

func TestCallExpr(t *testing.T) {
	var (
		data    = MustReadFile(t, "callExpr.go.txt")
		fileSet = token.NewFileSet()
		f       = MustParse(t, fileSet, "", data)
	)

	var (
		s1 = NodeSelections{}
		s2 = NodeSelections{}
		s3 = NodeSelections{}
	)

	Walk(f, PatternMatchers{
		New(
			[]NodeCondition{
				Type(new(ast.ExprStmt)),
				CallExpr(
					SelectorExpr(
						s1.Select(Ident("logrus"), "package1"),
						s1.Select(Ident("SetLevel"), "methodName"),
					),
					Exprs(
						[]NodeCondition{
							SelectorExpr(
								s1.Select(Ident("logrus"), "package2"),
								Ident("DebugLevel"),
							),
						},
					),
				),
			},
			func() {
				s1.Ident("package1").Name = "zerolog"
				s1.Ident("methodName").Name = "SetGlobalLevel"
				s1.Ident("package2").Name = "zerolog"
			},
		),
		New(
			[]NodeCondition{
				s2.Select(Type(new(ast.ExprStmt)), "call"),
				CallExpr(
					SelectorExpr(
						Ident("logrus"),
						s2.Select(IgnoreNode(), "methodName"),
					),
					Exprs(
						[]NodeCondition{
							s2.Select(IgnoreNode(), "arg"),
						},
					),
				),
			},
			func() {
				s2.ExprStmt("call").X = createZerologCallExpr(
					s2.Ident("methodName").Name,
					"Msg",
					s2.BasicLit("arg"),
				)
			},
		),
		New(
			[]NodeCondition{
				s3.Select(Type(new(ast.ExprStmt)), "call"),
				CallExpr(
					SelectorExpr(
						Ident("logrus"),
						s3.Select(Ident("Info"), "arg"),
					),
					Exprs(
						[]NodeCondition{
							s3.Select(IgnoreNode(), "arg1"),
							s3.Select(IgnoreNode(), "arg2"),
						},
					),
				),
			},
			func() {
				s3.ExprStmt("call").X = createZerologCallExpr(
					"Info",
					"Msgf",
					&ast.BasicLit{Kind: token.STRING, Value: `"%v %v"`},
					s3.BasicLit("arg1"),
					s3.BasicLit("arg2"))
			},
		),
	})

	patched := bytes.NewBuffer([]byte{})
	FailOnError(t, printer.Fprint(patched, fileSet, f))

	want := GetFunctionBody(t, data, "Want")
	got := GetFunctionBody(t, patched.Bytes(), "Got")

	AssertEquals(t, want, got)
}

func createZerologCallExpr(level string, method string, args ...ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X: &ast.CallExpr{
				Fun: &ast.SelectorExpr{
					X: &ast.Ident{
						Name: "log",
					},
					Sel: &ast.Ident{
						Name: level,
					},
				},
			},
			Sel: &ast.Ident{
				Name: method,
			},
		},
		Args: args,
	}
}
