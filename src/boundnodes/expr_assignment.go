package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Assignment expression
// ---------------------
type BoundAssignmentExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Expression BoundExpressionNode
    Value BoundExpressionNode
}

func NewBoundAssignmentExpressionNode(src syntaxnodes.SyntaxNode, expr BoundExpressionNode, val BoundExpressionNode) *BoundAssignmentExpressionNode {
    return &BoundAssignmentExpressionNode {
        SourceNode: src,
        Expression: expr,
        Value: val,
    }
}

func (nd *BoundAssignmentExpressionNode) Type() BoundNodeType {
    return BT_AssignmentExpr
}

func (nd *BoundAssignmentExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundAssignmentExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Value.ExprType()
} 
