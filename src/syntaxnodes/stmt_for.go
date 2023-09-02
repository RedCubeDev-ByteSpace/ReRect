package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ForStatementNode struct {
    StatementNode

    ForKw lexer.Token
    Declaration StatementNode
    Condition ExpressionNode
    Action StatementNode
    Body StatementNode
}

func NewForStatementNode(forkw lexer.Token, decl StatementNode, cond ExpressionNode, act StatementNode, body StatementNode) *ForStatementNode {
    return &ForStatementNode{
        ForKw: forkw,
        Declaration: decl,
        Condition: cond,
        Action: act,
        Body: body,
    }
}

func (n *ForStatementNode) Position() span.Span {
    return n.ForKw.Position.SpanBetween(n.Body.Position())
}

func (n *ForStatementNode) Type() SyntaxNodeType {
    return NT_ForStmt
}
