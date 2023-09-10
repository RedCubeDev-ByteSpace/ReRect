package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type FieldAssignmentClauseNode struct {
    SyntaxNode

    FieldName lexer.Token
    Value ExpressionNode
}

func NewFieldAssignmentClauseNode(fld lexer.Token, val ExpressionNode) *FieldAssignmentClauseNode {
    return &FieldAssignmentClauseNode{
        FieldName: fld,
        Value: val,
    }
}

func (n *FieldAssignmentClauseNode) Position() span.Span {
    return n.FieldName.Position.SpanBetween(n.Value.Position())
}

func (n *FieldAssignmentClauseNode) Type() SyntaxNodeType {
    return NT_FieldAssignmentCls
}

