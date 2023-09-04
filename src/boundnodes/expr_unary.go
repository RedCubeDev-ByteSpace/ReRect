package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Unary expression
// ----------------
type BoundUnaryExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Operator *BoundUnaryOperator
    Operand BoundExpressionNode
}

func NewBoundUnaryExpressionNode(src syntaxnodes.SyntaxNode, op *BoundUnaryOperator, operand BoundExpressionNode) *BoundUnaryExpressionNode {
    return &BoundUnaryExpressionNode {
        SourceNode: src,
        Operator: op,
        Operand: operand,
    }
}

func (nd *BoundUnaryExpressionNode) Type() BoundNodeType {
    return BT_UnaryExpr
}

func (nd *BoundUnaryExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundUnaryExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Operator.Result
} 
