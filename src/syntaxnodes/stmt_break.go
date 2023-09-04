package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type BreakStatementNode struct {
    ExpressionNode

    BreakKw lexer.Token
}

func NewBreakStatementNode(kw lexer.Token) *BreakStatementNode {
    return &BreakStatementNode{
        BreakKw: kw,
    }
}

func (n *BreakStatementNode) Position() span.Span {
    return n.BreakKw.Position
}

func (n *BreakStatementNode) Type() SyntaxNodeType {
    return NT_BreakStmt
}
