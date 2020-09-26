package asterisk

import (
	"go/ast"
	"reflect"
)

type (
	BoolCondition       func(bool) bool
	ChanDirCondition    func(ast.ChanDir) bool
	ExprCondition       func(ast.Expr) bool
	FilesMapCondition   func(Files map[string]*ast.File) bool
	ImportsMapCondition func(map[string]*ast.Object) bool
	NodeCondition       func(ast.Node) bool
	NodesCondition      func([]ast.Node) bool
	ScopeCondition      func(*ast.Scope) bool
	StringCondition     func(string) bool
)

/**************************************************************************
	concrete expression nodes
**************************************************************************/

// BadExpr check if the given ast.Node is a ast.BadExpr.
func BadExpr() NodeCondition {
	return Type(new(ast.BadExpr))
}

// Ident check if the given ast.Ident name matches the requested one.
func Ident(name string) NodeCondition {
	return func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			return ident.Name == name
		}

		return false
	}
}

// IdentExpr check if the given ast.IdentExpr name matches the requested one.
func IdentExpr(name string) NodeCondition {
	return func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			return ident.Name == name
		}

		return false
	}
}

// Ellipsis check if the given ast.Ellipsis matches the given conditions.
func Ellipsis(elem NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.Ellipsis); ok {
			return elem(e.Elt)
		}

		return false
	}
}

// BasicLit check if the given ast.BasicLit matches the given conditions.
func BasicLit(value string) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.BasicLit); ok {
			return e.Value == value
		}

		return false
	}
}

// FuncLit check if the given ast.FuncLit matches the given conditions.
func FuncLit(t, block NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.FuncLit); ok {
			return t(e.Type) && block(e.Body)
		}

		return false
	}
}

// CompositeLit check if the given ast.ParenExpr matches the given conditions.
func CompositeLit(t NodeCondition, args NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.CompositeLit); ok {
			return t(e.Type) && args(toNodes(e.Elts))
		}

		return false
	}
}

// ParenExpr check if the given ast.ParenExpr matches the given conditions.
func ParenExpr(x NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ParenExpr); ok {
			return x(e.X)
		}

		return false
	}
}

// Expr check if the given ast.Expr matches the given condition.
func Expr(x ExprCondition) ExprCondition {
	return func(n ast.Expr) bool {
		return x(n)
	}
}

// Exprs check if the given []ast.Node matches the given conditions in sequence.
func Exprs(x []NodeCondition) NodesCondition {
	return func(n []ast.Node) bool {
		if len(n) != len(x) {
			return false
		}

		for i := range n {
			if len(x) > i {
				if !x[i](n[i]) {
					return false
				}
			}
		}

		return true
	}
}

// SelectorExpr check if the given ast.SelectorExpr matches the given conditions.
func SelectorExpr(x, sel NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.SelectorExpr); ok {
			return x(e.X) && sel(e.Sel)
		}

		return false
	}
}

// IndexExpr check if the given ast.IndexExpr matches the given conditions.
func IndexExpr(x, index NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.IndexExpr); ok {
			return x(e.X) && index(e.Index)
		}

		return false
	}
}

// SliceExpr check if the given ast.SliceExpr matches the given conditions.
func SliceExpr(x, low, high, max NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.SliceExpr); ok {
			return x(e.X) && low(e.Low) && high(e.High) && max(e.Max)
		}

		return false
	}
}

// TypeAssertExpr check if the given ast.TypeAssertExpr matches the given conditions.
func TypeAssertExpr(x, t NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.TypeAssertExpr); ok {
			return x(e.X) && t(e.Type)
		}

		return false
	}
}

// CallExpr check if the given ast.CallExpr matches the given conditions.
func CallExpr(fun NodeCondition, args NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.CallExpr); ok {
			return fun(e.Fun) && args(toNodes(e.Args))
		}

		return false
	}
}

// StarExpr check if the given ast.TypeAssertExpr matches the given conditions.
func StarExpr(x NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.StarExpr); ok {
			return x(e.X)
		}

		return false
	}
}

// UnaryExpr check if the given ast.UnaryExpr matches the given conditions.
func UnaryExpr(x NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.UnaryExpr); ok {
			return x(e.X)
		}

		return false
	}
}

// BinaryExpr check if the given ast.BinaryExpr matches the given conditions.
func BinaryExpr(x, y NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.BinaryExpr); ok {
			return x(e.X) && y(e.Y)
		}

		return false
	}
}

// KeyValueExpr check if the given ast.KeyValueExpr matches the given conditions.
func KeyValueExpr(k, v NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.KeyValueExpr); ok {
			return k(e.Key) && v(e.Value)
		}

		return false
	}
}

/**************************************************************************
	type-specific expression nodes
**************************************************************************/

// ArrayType check if the given ast.ArrayType matches the given conditions.
func ArrayType(elt, l NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ArrayType); ok {
			return elt(e.Elt) && l(e.Len)
		}

		return false
	}
}

// StructType check if the given ast.StructType matches the given conditions.
func StructType(fields NodeCondition, incomplete BoolCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.StructType); ok {
			return fields(e.Fields) && incomplete(e.Incomplete)
		}

		return false
	}
}

// FuncType check if the given ast.FuncType matches the given conditions.
func FuncType(params, results NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.FuncType); ok {
			return params(e.Params) && results(e.Results)
		}

		return false
	}
}

// InterfaceType check if the given ast.InterfaceType matches the given conditions.
func InterfaceType(methods NodeCondition, incomplete BoolCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.InterfaceType); ok {
			return methods(e.Methods) && incomplete(e.Incomplete)
		}

		return false
	}
}

// MapType check if the given ast.MapType matches the given conditions.
func MapType(k, v NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.MapType); ok {
			return k(e.Key) && v(e.Value)
		}

		return false
	}
}

// ChanType check if the given ast.ChanType matches the given conditions.
func ChanType(k NodeCondition, v ChanDirCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ChanType); ok {
			return k(e.Value) && v(e.Dir)
		}

		return false
	}
}

/**************************************************************************
	concrete statement nodes
**************************************************************************/

// BadStmt check if the given ast.Node is a ast.BadStmt.
func BadStmt() NodeCondition {
	return Type(new(ast.BadStmt))
}

// DeclStmt check if the given ast.DeclStmt matches the given conditions.
func DeclStmt(decl NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.DeclStmt); ok {
			return decl(e.Decl)
		}

		return false
	}
}

// EmptyStmt check if the given ast.EmptyStmt matches the given conditions.
func EmptyStmt(implicit BoolCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.EmptyStmt); ok {
			return implicit(e.Implicit)
		}

		return false
	}
}

// LabeledStmt check if the given ast.LabeledStmt matches the given conditions.
func LabeledStmt(label, stmt NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.LabeledStmt); ok {
			return label(e.Label) && stmt(e.Label)
		}

		return false
	}
}

// ExprStmt check if the given ast.ExprStmt matches the given conditions.
func ExprStmt(x NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ExprStmt); ok {
			return x(e.X)
		}

		return false
	}
}

// SendStmt check if the given ast.SendStmt matches the given conditions.
func SendStmt(channel, val NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.SendStmt); ok {
			return channel(e.Chan) && val(e.Value)
		}

		return false
	}
}

// IncDecStmt check if the given ast.IncDecStmt matches the given conditions.
func IncDecStmt(x NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.IncDecStmt); ok {
			return x(e.X)
		}

		return false
	}
}

// AssignStmt check if the given ast.AssignStmt matches the given conditions.
func AssignStmt(lhs, rhs NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.AssignStmt); ok {
			return lhs(toNodes(e.Lhs)) && rhs(toNodes(e.Rhs))
		}

		return false
	}
}

// GoStmt check if the given ast.GoStmt matches the given conditions.
func GoStmt(call NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.GoStmt); ok {
			return call(e.Call)
		}

		return false
	}
}

// DeferStmt check if the given ast.DeferStmt matches the given conditions.
func DeferStmt(call NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.DeferStmt); ok {
			return call(e.Call)
		}

		return false
	}
}

// ReturnStmt check if the given ast.ReturnStmt matches the given conditions.
func ReturnStmt(results NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ReturnStmt); ok {
			return results(toNodes(e.Results))
		}

		return false
	}
}

// BranchStmt check if the given ast.BranchStmt matches the given conditions.
func BranchStmt(label NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.BranchStmt); ok {
			return label(e.Label)
		}

		return false
	}
}

// BlockStmt check if the given ast.BranchStmt matches the given conditions.
func BlockStmt(stmts NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.BlockStmt); ok {
			return stmts(toNodes(e.List))
		}

		return false
	}
}

// IfStmt check if the given ast.IfStmt matches the given conditions.
func IfStmt(init, body, cond, els NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.IfStmt); ok {
			return init(e.Init) && cond(e.Cond) && body(e.Body) && els(e.Else)
		}

		return false
	}
}

// CaseClause check if the given ast.CaseClause matches the given conditions.
func CaseClause(list NodesCondition, body NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.CaseClause); ok {
			return list(toNodes(e.List)) && body(toNodes(e.Body))
		}

		return false
	}
}

// SwitchStmt check if the given ast.SwitchStmt matches the given conditions.
func SwitchStmt(init, tag, body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.SwitchStmt); ok {
			return init(e.Init) && tag(e.Tag) && body(e.Body)
		}

		return false
	}
}

// TypeSwitchStmt check if the given ast.SwitchStmt matches the given conditions.
func TypeSwitchStmt(init, assign, body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.TypeSwitchStmt); ok {
			return init(e.Init) && assign(e.Assign) && body(e.Body)
		}

		return false
	}
}

// CommClause check if the given ast.CommClause matches the given conditions.
func CommClause(comm NodeCondition, body NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.CommClause); ok {
			return comm(e.Comm) && body(toNodes(e.Body))
		}

		return false
	}
}

// SelectStmt check if the given ast.SelectStmt matches the given conditions.
func SelectStmt(body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.SelectStmt); ok {
			return body(e.Body)
		}

		return false
	}
}

// ForStmt check if the given ast.ForStmt matches the given conditions.
func ForStmt(init, cond, post, body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ForStmt); ok {
			return init(e.Init) && cond(e.Cond) && post(e.Post) && body(e.Body)
		}

		return false
	}
}

// RangeStmt check if the given ast.RangeStmt matches the given conditions.
func RangeStmt(k, v, x, body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.RangeStmt); ok {
			return k(e.Key) && v(e.Value) && x(e.X) && body(e.Body)
		}

		return false
	}
}

/**************************************************************************
 single (non-parenthesized) import, constant, type, or variable declaration
**************************************************************************/

// Spec check if the given values type matches the requested one.
func Spec() NodeCondition {
	return Type(new(ast.Spec))
}

// ImportSpec check if the given ast.ImportSpec matches the given conditions.
func ImportSpec(doc, name, importPath, comment NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ImportSpec); ok {
			return doc(e.Doc) && name(e.Name) && importPath(e.Path) && comment(e.Comment)
		}

		return false
	}
}

// ValueSpec check if the given ast.ValueSpec matches the given conditions.
func ValueSpec(doc, t, comment NodeCondition, names NodesCondition, values NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.ValueSpec); ok {
			return doc(e.Doc) && names(toNodes(e.Names)) && t(e.Type) && values(toNodes(e.Values)) && comment(e.Comment)
		}

		return false
	}
}

// TypeSpec check if the given ast.TypeSpec matches the given conditions.
func TypeSpec(doc, name, t, comment NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.TypeSpec); ok {
			return doc(e.Doc) && name(e.Name) && t(e.Type) && comment(e.Comment)
		}

		return false
	}
}

/**************************************************************************
	declaration nodes
**************************************************************************/

// BadDecl check if the given values type matches *ast.BadDecl.
func BadDecl() NodeCondition {
	return Type(new(ast.BadDecl))
}

// GenDecl check if the given ast.GenDecl matches the given conditions.
func GenDecl(doc NodeCondition, specs NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.GenDecl); ok {
			return doc(e.Doc) && specs(toNodes(e.Specs))
		}

		return false
	}
}

// FuncDecl check if the given ast.FuncDecl matches the given conditions.
func FuncDecl(doc, recv, name, t, body NodeCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.FuncDecl); ok {
			return doc(e.Doc) && recv(e.Recv) && name(e.Name) && t(e.Type) && body(e.Body)
		}

		return false
	}
}

/**************************************************************************
	Files and packages
**************************************************************************/

// File check if the given ast.File matches the given conditions.
func File(
	doc,
	name NodeCondition,
	decls NodesCondition,
	scope ScopeCondition,
	imports,
	unresolved NodesCondition,
	comments NodesCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.File); ok {
			return doc(e.Doc) &&
				name(e.Name) &&
				decls(toNodes(e.Decls)) &&
				scope(e.Scope) &&
				imports(toNodes(e.Imports)) &&
				unresolved(toNodes(e.Unresolved)) &&
				comments(toNodes(e.Comments))
		}

		return false
	}
}

// Package check if the given ast.Package matches the given conditions.
func Package(
	scope ScopeCondition,
	name StringCondition,
	imports ImportsMapCondition,
	files FilesMapCondition) NodeCondition {
	return func(n ast.Node) bool {
		if e, ok := n.(*ast.Package); ok {
			return scope(e.Scope) && name(e.Name) && imports(e.Imports) && files(e.Files)
		}

		return false
	}
}

/**************************************************************************
	custom
**************************************************************************/

func First(c NodeCondition) NodesCondition {
	return func(nodes []ast.Node) bool {
		if len(nodes) == 0 {
			return true
		}

		return c(nodes[0])
	}
}

func Last(c NodeCondition) NodesCondition {
	return func(nodes []ast.Node) bool {
		if len(nodes) == 0 {
			return true
		}

		return c(nodes[len(nodes)-1])
	}
}

// Type check if the given values type matches the requested one.
func Type(t interface{}) NodeCondition {
	wantType := reflect.ValueOf(t).Type()

	return func(n ast.Node) bool {
		return wantType.AssignableTo(reflect.ValueOf(n).Type())
	}
}

// IgnoreNode always returns true.
func IgnoreNode() NodeCondition {
	return func(n ast.Node) bool {
		return true
	}
}

func IgnoreNodes() NodesCondition {
	return func(n []ast.Node) bool {
		return true
	}
}

func IgnoreScope() ScopeCondition {
	return func(n *ast.Scope) bool {
		return true
	}
}

func toNodes(e interface{}) []ast.Node {
	var (
		n  []ast.Node
		ev = reflect.ValueOf(e)
	)

	for i := 0; i < ev.Len(); i++ {
		switch v := ev.Index(i).Interface().(type) {
		case *ast.Ident:
			n = append(n, v)
		case *ast.ImportSpec:
			n = append(n, v)
		case *ast.CommentGroup:
			n = append(n, v)
		case ast.Expr:
			n = append(n, v)
		case ast.Decl:
			n = append(n, v)
		case ast.Stmt:
			n = append(n, v)
		}
	}

	return n
}
