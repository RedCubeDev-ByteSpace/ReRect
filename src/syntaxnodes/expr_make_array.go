package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type MakeArrayExpressionNode struct {
    ExpressionNode

    MakeKw lexer.Token
    ClosingTok lexer.Token

    ArrType *TypeClauseNode
    
    Length ExpressionNode
    Initializers []ExpressionNode
    HasInitializers bool
}

func NewMakeArrayExpressionNode(makekw lexer.Token, cls lexer.Token, typ *TypeClauseNode, length ExpressionNode, init []ExpressionNode, hasinit bool) *MakeArrayExpressionNode {
    return &MakeArrayExpressionNode{
        MakeKw: makekw,
        ClosingTok: cls,
        ArrType: typ,
        Length: length,
        Initializers: init,
        HasInitializers: hasinit,
    }
}

func (n *MakeArrayExpressionNode) Position() span.Span {
    return n.MakeKw.Position.SpanBetween(n.ClosingTok.Position)
}

func (n *MakeArrayExpressionNode) Type() SyntaxNodeType {
    return NT_MakeArrayExpr
}
