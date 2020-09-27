package test

import (
	"bytes"
	"go/printer"
	"go/token"
	"testing"

	. "github.com/Oppodelldog/asterisk"
)

func TestImports(t *testing.T) {
	var (
		file     = "code.go"
		data     = MustReadFile(t, "imports1.go.txt")
		wantData = MustReadFile(t, "imports2.go.txt")
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
