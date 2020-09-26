package callExpr

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	. "github.com/Oppodelldog/asterisk"
	. "github.com/Oppodelldog/asterisk/test"
)

func TestAsterisk(t *testing.T) {
	var (
		file    = "code.go"
		data    = MustReadFile(t, "code.go")
		fileSet = token.NewFileSet()
		f       = MustParse(t, fileSet, file, data)
		s1      = NodeSelections{}
	)

	Walk(f, PatternMatchers{
		New(
			[]NodeCondition{
				FuncDecl(
					IgnoreNode(),
					IgnoreNode(),
					IgnoreNode(),
					IgnoreNode(),
					s1.Select(
						BlockStmt(
							First(
								s1.Select(
									IfStmt(
										IgnoreNode(),
										s1.Select(IgnoreNode(), "body"),
										IgnoreNode(),
										s1.Select(
											BlockStmt(
												Last(
													ReturnStmt(IgnoreNodes()),
												),
											),
											"else",
										),
									),
									"if",
								),
							),
						),
						"block",
					),
				),
			},
			func() {
				var elseStmts = s1.BlockStmt("else").List
				var ret = elseStmts[len(elseStmts)-1]
				s1.BlockStmt("block").List = append(s1.BlockStmt("block").List, ret)
				s1.IfStmt("if").Else = nil
			},
		),
	})

	patched := bytes.NewBuffer([]byte{})
	FailOnError(t, printer.Fprint(patched, fileSet, f))

	want := GetFunctionBody(t, data, "Want")
	got := GetFunctionBody(t, patched.Bytes(), "Got")

	AssertEquals(t, want, got)
}
