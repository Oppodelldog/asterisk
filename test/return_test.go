package test

import (
	"go/token"
	"testing"

	. "github.com/Oppodelldog/asterisk"
)

func TestReturnStmt_detectEmptyLines(t *testing.T) {
	var (
		data    = MustReadFile(t, "return.go.txt")
		fileSet = token.NewFileSet()
		f       = MustParse(t, fileSet, "", data)
		s1      = NodeSelections{}
	)

	Walk(f, PatternMatchers{
		New(
			[]NodeCondition{
				s1.Select(
					BlockStmt(
						Last(ReturnStmt(IgnoreNodes())),
					), "blockWithReturn",
				),
			},
			func() {
				var (
					block        = s1.BlockStmt("blockWithReturn")
					stmts        = block.List
					last         = stmts[0]
					lineAt       = fileSet.Position(last.End()).Line
					lineEndBlock = fileSet.Position(block.Rbrace).Line
					lineBefore   int
				)

				if len(stmts) == 1 {
					lineBefore = fileSet.Position(block.Lbrace).Line
				} else {
					beforeLast := stmts[len(stmts)-2]
					lineBefore = fileSet.Position(beforeLast.End()).Line
				}

				if lineAt-lineBefore > 1 {
					t.Logf("unnecessary empty line before return (%v - %v)", lineAt, lineEndBlock)
				}
				if lineEndBlock-lineAt > 1 {
					t.Logf("unnecessary empty line before closing brace (%v - %v)", lineAt, lineEndBlock)
				}
			},
		),
	})
}
