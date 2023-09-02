package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FromToStatementNode struct {
    StatementNode

    FromKw lexer.Token
    LowerBound ExpressionNode
    UpperBound ExpressionNode
    Body StatementNode
}

func NewFromToStatementNode(fromkw lexer.Token, lwbound ExpressionNode, upbound ExpressionNode, body StatementNode) FromToStatementNode {
    return FromToStatementNode{
        FromKw: fromkw,
        LowerBound: lwbound,
        UpperBound: upbound,
        Body: body,
    }
}

func (n *FromToStatementNode) Position() span.Span {
    return n.FromKw.Position.SpanBetween(n.Body.Position())
}

func (n *FromToStatementNode) Type() SyntaxNodeType {
    return NT_FromToStmt
}
