package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type WhileStatementNode struct {
    StatementNode

    WhileKw lexer.Token
    Expression ExpressionNode
    Body StatementNode
}

func NewWhileStatementNode(whilekw lexer.Token, expr ExpressionNode, body StatementNode) WhileStatementNode {
    return WhileStatementNode{
        WhileKw: whilekw,
        Expression: expr,
        Body: body,
    }
}

func (n *WhileStatementNode) Position() span.Span {
    return n.WhileKw.Position.SpanBetween(n.Body.Position())
}

func (n *WhileStatementNode) Type() SyntaxNodeType {
    return NT_WhileStmt
}
