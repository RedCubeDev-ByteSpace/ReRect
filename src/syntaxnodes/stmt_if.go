package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type IfStatementNode struct {
    StatementNode

    IfKw lexer.Token
    Expression ExpressionNode
    Body StatementNode
    Else StatementNode
    HasElseClause bool
}

func NewIfStatementNode(ifkw lexer.Token, expr ExpressionNode, body StatementNode, els StatementNode, haselse bool) IfStatementNode {
    return IfStatementNode{
        IfKw: ifkw,
        Expression: expr,
        Body: body,
        Else: els,
        HasElseClause: haselse,
    }
}

func (n *IfStatementNode) Position() span.Span {
    spn := n.IfKw.Position.SpanBetween(n.Body.Position())

    if n.HasElseClause {
        spn = spn.SpanBetween(n.Else.Position())
    }

    return spn
}

func (n *IfStatementNode) Type() SyntaxNodeType {
    return NT_IfStmt
}
