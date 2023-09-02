package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type AssignmentExpressionNode struct {
    ExpressionNode

    VarName lexer.Token
    Expression ExpressionNode
}

func NewAssignmentExpressionNode(varname lexer.Token, expr ExpressionNode) AssignmentExpressionNode {
    return AssignmentExpressionNode{
        VarName: varname,
        Expression: expr,
    }
}

func (n *AssignmentExpressionNode) Position() span.Span {
    return n.VarName.Position.SpanBetween(n.Expression.Position())
}

func (n *AssignmentExpressionNode) Type() SyntaxNodeType {
    return NT_AssignmentExpr
}
