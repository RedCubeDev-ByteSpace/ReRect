package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ReturnStatementNode struct {
    StatementNode

    ReturnKw lexer.Token

    Expression ExpressionNode
    HasExpression bool
}

func NewReturnStatementNode(retkw lexer.Token, expr ExpressionNode, hasexpr bool) *ReturnStatementNode {
    return &ReturnStatementNode{
        ReturnKw: retkw,
        Expression: expr,
        HasExpression: hasexpr,
    }
}

func (n *ReturnStatementNode) Position() span.Span {
    spn := n.ReturnKw.Position

    if n.HasExpression {
        spn = spn.SpanBetween(n.Expression.Position())
    }

    return spn
}

func (n *ReturnStatementNode) Type() SyntaxNodeType {
    return NT_ReturnStmt
}
