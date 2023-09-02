package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type LoopStatementNode struct {
    StatementNode

    LoopKw lexer.Token
    Expression ExpressionNode
    Body StatementNode
}

func NewLoopStatementNode(loopkw lexer.Token, expr ExpressionNode, body StatementNode) *LoopStatementNode {
    return &LoopStatementNode{
        LoopKw: loopkw,
        Expression: expr,
        Body: body,
    }
}

func (n *LoopStatementNode) Position() span.Span {
    return n.LoopKw.Position.SpanBetween(n.Body.Position())
}

func (n *LoopStatementNode) Type() SyntaxNodeType {
    return NT_LoopStmt
}
