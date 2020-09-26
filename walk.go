package asterisk

import (
	"go/ast"
)

func Walk(f ast.Node, pms PatternMatchers) {
	ast.Inspect(f, func(n ast.Node) bool {
		if n == nil {
			return true
		}

		pms.Match(n)

		return true
	})
}
