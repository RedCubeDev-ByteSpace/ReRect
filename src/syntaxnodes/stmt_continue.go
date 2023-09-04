package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type ContinueStatementNode struct {
    ExpressionNode

    ContinueKw lexer.Token
}

func NewContinueStatementNode(kw lexer.Token) *ContinueStatementNode {
    return &ContinueStatementNode{
        ContinueKw: kw,
    }
}

func (n *ContinueStatementNode) Position() span.Span {
    return n.ContinueKw.Position
}

func (n *ContinueStatementNode) Type() SyntaxNodeType {
    return NT_ContinueStmt
}
