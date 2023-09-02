package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FromToStatementNode struct {
    StatementNode

    FromKw lexer.Token
    LowerBound ExpressionNode
    Iterator lexer.Token
    UpperBound ExpressionNode
    Body StatementNode
}

func NewFromToStatementNode(fromkw lexer.Token, lwbound ExpressionNode, iterator lexer.Token, upbound ExpressionNode, body StatementNode) *FromToStatementNode {
    return &FromToStatementNode{
        FromKw: fromkw,
        LowerBound: lwbound,
        Iterator: iterator,
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
