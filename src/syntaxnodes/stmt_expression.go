package syntaxnodes

import (
	"bytespace.network/rerect/span"
)

type ExpressionStatementNode struct {
    StatementNode

    Expression ExpressionNode
}

func NewExpressionStatementNode(expr ExpressionNode) ExpressionStatementNode {
    return ExpressionStatementNode{
        Expression: expr,
    }
}

func (n *ExpressionStatementNode) Position() span.Span {
    return n.Expression.Position()
}

func (n *ExpressionStatementNode) Type() SyntaxNodeType {
    return NT_ExpressionStmt
}
