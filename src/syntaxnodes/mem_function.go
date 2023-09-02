package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FunctionNode struct {
    MemberNode

    FunctionKw lexer.Token
    FunctionName lexer.Token

    Parameters []*ParameterClauseNode 

    ReturnType *TypeClauseNode
    HasReturnType bool

    Body StatementNode
}

func NewFunctionNode(fnckw lexer.Token, fncname lexer.Token, prm []*ParameterClauseNode, rettype *TypeClauseNode, hasrettype bool, body StatementNode) *FunctionNode {
    return &FunctionNode{
        FunctionKw: fnckw,
        FunctionName: fncname,
        Parameters: prm,
        ReturnType: rettype,
        HasReturnType: hasrettype,
        Body: body,
    }
}

func (n *FunctionNode) Position() span.Span {
    return n.FunctionKw.Position.SpanBetween(n.Body.Position())
}

func (n *FunctionNode) Type() SyntaxNodeType {
    return NT_Function
}
