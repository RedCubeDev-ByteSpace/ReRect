package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Binary expression
// -----------------
type BoundBinaryExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Operator BoundBinaryOperator
    Left BoundExpressionNode
    Right BoundExpressionNode
}

func NewBoundBinaryExpressionNode(src syntaxnodes.SyntaxNode, op BoundBinaryOperator, left BoundExpressionNode, right BoundExpressionNode) *BoundBinaryExpressionNode {
    return &BoundBinaryExpressionNode {
        SourceNode: src,
        Operator: op,
        Left: left,
        Right: right,
    }
}

func (nd *BoundBinaryExpressionNode) Type() BoundNodeType {
    return BT_BinaryExpr
}

func (nd *BoundBinaryExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundBinaryExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Operator.Result
} 
