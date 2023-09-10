package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type MakeExpressionNode struct {
    ExpressionNode

    MakeKw lexer.Token
    ClosingTok lexer.Token

    Package lexer.Token
    Container lexer.Token
    HasPackage bool

    Initializer []*FieldAssignmentClauseNode
    HasInitializer bool

    ConstructorArguments []ExpressionNode
    HasConstructor bool
}

func NewMakeExpressionNode(kw lexer.Token, cls lexer.Token, cnt lexer.Token, pck lexer.Token, haspck bool, init []*FieldAssignmentClauseNode, hasinit bool, args []ExpressionNode, hascst bool) *MakeExpressionNode {
    return &MakeExpressionNode{
        MakeKw: kw,
        ClosingTok: cls,
        Package: pck,
        Container: cnt,
        HasPackage: haspck,
        Initializer: init,
        HasInitializer: hasinit,
        ConstructorArguments: args,
        HasConstructor: hascst,
    }
}

func (n *MakeExpressionNode) Position() span.Span {
    return n.MakeKw.Position.SpanBetween(n.ClosingTok.Position)
}

func (n *MakeExpressionNode) Type() SyntaxNodeType {
    return NT_MakeExpr
}
