package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// ArrayIndex expression
// ------------------
type BoundArrayIndexExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    SourceArray BoundExpressionNode
    Index BoundExpressionNode
}

func NewBoundArrayIndexExpressionNode(src syntaxnodes.SyntaxNode, srcarr BoundExpressionNode, idx BoundExpressionNode) *BoundArrayIndexExpressionNode {
    return &BoundArrayIndexExpressionNode {
        SourceNode: src,
        SourceArray: srcarr,
        Index: idx,
    }
}

func (nd *BoundArrayIndexExpressionNode) Type() BoundNodeType {
    return BT_ArrayIndexExpr
}

func (nd *BoundArrayIndexExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundArrayIndexExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.SourceArray.ExprType().SubTypes[0]
} 
