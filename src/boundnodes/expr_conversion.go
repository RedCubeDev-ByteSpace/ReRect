package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Conversion expression
// ---------------------
type BoundConversionExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Value BoundExpressionNode
    TargetType *symbols.TypeSymbol
}

func NewBoundConversionExpressionNode(src syntaxnodes.SyntaxNode, val BoundExpressionNode, target *symbols.TypeSymbol) *BoundConversionExpressionNode {
    return &BoundConversionExpressionNode {
        SourceNode: src,
        Value: val,
        TargetType: target,
    }
}

func (nd *BoundConversionExpressionNode) Type() BoundNodeType {
    return BT_ConversionExpr
}

func (nd *BoundConversionExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundConversionExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.TargetType
} 
