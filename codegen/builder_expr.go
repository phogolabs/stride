package codegen

import (
	"go/token"

	"github.com/dave/dst"
)

// Expr returns a expression
type Expr interface {
	Expr() dst.Expr
}

// PairExpr represents a pair expression
type PairExpr struct {
	X Expr
	Y Expr
}

// Pair returns the pair
func Pair(x Expr, y Expr) *PairExpr {
	return &PairExpr{
		X: x,
		Y: y,
	}
}

// Map represents a map expression
type Map map[string]Expr

var _ Expr = LeafExpr("")

// LeafExpr represents node expression
type LeafExpr string

// Leaf returns a text expr
func Leaf(text string) LeafExpr {
	return LeafExpr(text)
}

// Expr return the expression
func (n LeafExpr) Expr() dst.Expr {
	return &dst.Ident{
		Name: string(n),
	}
}

// Stmt return the statement
func (n LeafExpr) Stmt() dst.Stmt {
	return &dst.ExprStmt{
		X: &dst.Ident{
			Name: string(n),
		},
	}
}

var (
	_ Expr = &CallExpr{}
	_ Stmt = &CallExpr{}
)

// CallExpr represents a call expression
type CallExpr struct {
	node *dst.CallExpr
}

// Expr return the expresion
func (call *CallExpr) Expr() dst.Expr {
	return call.node
}

// Stmt return the statement
func (call *CallExpr) Stmt() dst.Stmt {
	return &dst.ExprStmt{
		X: call.node,
	}
}

// Call returns a call expression
func Call(method Expr, args ...Expr) *CallExpr {
	expr := &CallExpr{
		node: &dst.CallExpr{
			Fun: method.Expr(),
		},
	}

	for _, param := range args {
		expr.node.Args = append(expr.node.Args, param.Expr())
	}

	return expr
}

var _ Expr = &BinaryExpr{}

// BinaryExpr represent a binary expression
type BinaryExpr struct {
	node *dst.BinaryExpr
}

// Condition returns the binary expr
func Condition(x Expr, op token.Token, y Expr) *BinaryExpr {
	return &BinaryExpr{
		node: &dst.BinaryExpr{
			X:  x.Expr(),
			Op: op,
			Y:  y.Expr(),
		},
	}
}

// Expr returns the expression
func (expr *BinaryExpr) Expr() dst.Expr {
	return expr.node
}
