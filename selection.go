package asterisk

import "go/ast"

// NodeSelections contain nodes that were selected during matching.
// They are stored as Pointers on the existing ast.Node Pointers so they can be used to
// manipulate the tree in memory.
type NodeSelections map[string][]**ast.Node

// BasicLit returns a pointer to the ast.Basic that was selected using the given key.
func (s NodeSelections) BasicLit(key string) *ast.BasicLit {
	return (**s[key][0]).(*ast.BasicLit)
}

// Ident returns a pointer to the ast.Ident that was selected using the given key.
func (s NodeSelections) Ident(key string) *ast.Ident {
	return (**s[key][0]).(*ast.Ident)
}

// CallExpr returns a pointer to the ast.CallExpr that was selected using the given key.
func (s NodeSelections) CallExpr(key string) *ast.CallExpr {
	return (**s[key][0]).(*ast.CallExpr)
}

// ExprStmt returns a pointer to the ast.ExprStmt that was selected using the given key.
func (s NodeSelections) ExprStmt(key string) *ast.ExprStmt {
	return (**s[key][0]).(*ast.ExprStmt)
}

// BlockStmt returns a pointer to the ast.BlockStmt that was selected using the given key.
func (s NodeSelections) BlockStmt(key string) *ast.BlockStmt {
	return (**s[key][0]).(*ast.BlockStmt)
}

// Stmt returns a pointer to the ast.Stmt that was selected using the given key.
func (s NodeSelections) Stmt(key string) ast.Stmt {
	return (**s[key][0]).(ast.Stmt)
}

// IfStmt returns a pointer to the ast.IfStmt that was selected using the given key.
func (s NodeSelections) IfStmt(key string) *ast.IfStmt {
	return (**s[key][0]).(*ast.IfStmt)
}

// Select will select the visited node for the given key if the given condition matches.
func (s NodeSelections) Select(c NodeCondition, key string) NodeCondition {
	return func(n ast.Node) bool {
		var res = c(n)

		if res {
			var (
				nodes []**ast.Node
				n1    = &n
			)

			nodes = append(nodes, &n1)
			s[key] = nodes
		}

		return res
	}
}

// ExprStmt returns a pointer to the ast.ExprStmt that was selected using the given key.
func (s NodeSelections) ImportSpecs(key string) []*ast.ImportSpec {
	var (
		nodes       = s[key]
		importSpecs = make([]*ast.ImportSpec, len(nodes))
	)

	for i := 0; i < len(nodes); i++ {
		importSpecs[i] = (**nodes[i]).(*ast.ImportSpec)
	}

	return importSpecs
}

// Select will select the visited nodes for the given key if the given condition matches.
func (s NodeSelections) Selects(c NodesCondition, key string) NodesCondition {
	return func(n []ast.Node) bool {
		var res = c(n)

		if res {
			var (
				nodes []**ast.Node
			)

			for i := range n {
				var n1 = &n[i]
				nodes = append(nodes, &n1)
			}

			s[key] = nodes
		}

		return res
	}
}
