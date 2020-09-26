package callExpr

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	. "github.com/Oppodelldog/asterisk"
	. "github.com/Oppodelldog/asterisk/test"
)

func TestChangeImportName(t *testing.T) {
	var (
		file     = "code.go"
		data     = MustReadFile(t, "code1.go")
		wantData = MustReadFile(t, "code2.go")
		fileSet  = token.NewFileSet()
		f        = MustParse(t, fileSet, file, data)
		s1       = NodeSelections{}
	)

	Walk(f, PatternMatchers{
		New(
			[]NodeCondition{
				File(
					IgnoreNode(),
					IgnoreNode(),
					IgnoreNodes(),
					IgnoreScope(),
					s1.Selects(IgnoreNodes(), "imports"),
					IgnoreNodes(),
					IgnoreNodes(),
				),
			},
			func() {
				s1.ImportSpecs("imports")[0].Name.Name = "changedImportName"
				s1.ImportSpecs("imports")[0].Path.Value = `"fmt"`
			},
		),
	})

	patched := bytes.NewBuffer([]byte{})
	FailOnError(t, printer.Fprint(patched, fileSet, f))

	want := GetImports(t, wantData)
	got := GetImports(t, patched.Bytes())

	AssertEquals(t, want, got)
}
