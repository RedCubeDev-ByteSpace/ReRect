package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FunctionNode struct {
    MemberNode

    FunctionKw lexer.Token
    FunctionName lexer.Token
    IsConstructor bool

    Parameters []*ParameterClauseNode 

    ReturnType *TypeClauseNode
    HasReturnType bool

    Body StatementNode
    HasBody bool
    Closing lexer.Token
}

func NewFunctionNode(fnckw lexer.Token, fncname lexer.Token, iscst bool, prm []*ParameterClauseNode, rettype *TypeClauseNode, hasrettype bool, body StatementNode, hasbody bool, closing lexer.Token) *FunctionNode {
    return &FunctionNode{
        FunctionKw: fnckw,
        FunctionName: fncname,
        IsConstructor: iscst,
        Parameters: prm,
        ReturnType: rettype,
        HasReturnType: hasrettype,
        Body: body,
        HasBody: hasbody,
        Closing: closing,
    }
}

func (n *FunctionNode) Position() span.Span {
    if n.HasBody {
        return n.FunctionKw.Position.SpanBetween(n.Body.Position())
    } else {
        return n.FunctionKw.Position.SpanBetween(n.Closing.Position)
    }
}

func (n *FunctionNode) Type() SyntaxNodeType {
    return NT_Function
}
