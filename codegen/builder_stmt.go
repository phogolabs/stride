package codegen

import (
	"go/token"

	"github.com/dave/dst"
)

// Stmt returns a statemtn
type Stmt interface {
	Stmt() dst.Stmt
}

var _ Stmt = &IfStmt{}

// IfStmt represents a if expression
type IfStmt struct {
	node *dst.IfStmt
}

// If returns an if expresssion
func If(cond Expr) *IfStmt {
	stmt := &IfStmt{
		node: &dst.IfStmt{
			Cond: cond.Expr(),
		},
	}

	stmt.node.Decs.Before = dst.EmptyLine
	stmt.node.Decs.After = dst.EmptyLine
	return stmt
}

// Stmt returns the statement
func (iff *IfStmt) Stmt() dst.Stmt {
	return iff.node
}

// Init init var
func (iff *IfStmt) Init(stmt Stmt) *IfStmt {
	iff.node.Init = stmt.Stmt()
	return iff
}

// Then inits the block
func (iff *IfStmt) Then(stmts ...Stmt) *IfStmt {
	block := &dst.BlockStmt{}
	iff.node.Body = block

	for _, statement := range stmts {
		block.List = append(block.List, statement.Stmt())
	}

	return iff
}

// Else inits the else block
func (iff *IfStmt) Else(stmts ...Stmt) *IfStmt {
	block := &dst.BlockStmt{}
	iff.node.Else = block

	for _, statement := range stmts {
		block.List = append(block.List, statement.Stmt())
	}

	return iff
}

var _ Stmt = &AssignStmt{}

// AssignStmt represents an assign expression
type AssignStmt struct {
	node *dst.AssignStmt
}

// Assign returns an assign expression
func Assign(expr ...*PairExpr) *AssignStmt {
	stmt := &AssignStmt{
		node: &dst.AssignStmt{
			Tok: token.DEFINE,
			Lhs: []dst.Expr{},
			Rhs: []dst.Expr{},
		},
	}

	for _, pair := range expr {
		stmt.node.Lhs = append(stmt.node.Lhs, pair.X.Expr())
		stmt.node.Rhs = append(stmt.node.Rhs, pair.Y.Expr())
	}

	return stmt
}

// Stmt returns the statement
func (assign *AssignStmt) Stmt() dst.Stmt {
	return assign.node
}

var _ Stmt = &DeclareStmt{}

// DeclareStmt represents a declare expression
type DeclareStmt struct {
	node *dst.DeclStmt
}

// Declare returns an declare expression
func Declare(expr Map) *DeclareStmt {
	specs := []dst.Spec{}

	for k, v := range expr {
		spec := &dst.ValueSpec{
			Names: []*dst.Ident{
				&dst.Ident{
					Name: k,
				},
			},
			Values: []dst.Expr{
				v.Expr(),
			},
		}

		specs = append(specs, spec)
	}

	stmt := &DeclareStmt{
		node: &dst.DeclStmt{
			Decl: &dst.GenDecl{
				Tok:   token.VAR,
				Specs: specs,
			},
		},
	}

	stmt.node.Decs.Before = dst.EmptyLine
	stmt.node.Decs.After = dst.EmptyLine

	return stmt
}

// Stmt returns the statement
func (assign *DeclareStmt) Stmt() dst.Stmt {
	return assign.node
}
