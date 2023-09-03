package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Literal expression
// ------------------
type BoundLiteralExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    LiteralType *symbols.TypeSymbol
    LiteralValue interface{}
}

func NewBoundLiteralExpressionNode(src syntaxnodes.SyntaxNode, littype *symbols.TypeSymbol, litv interface{}) *BoundLiteralExpressionNode {
    return &BoundLiteralExpressionNode {
        SourceNode: src,
        LiteralType: littype,
        LiteralValue: litv,
    }
}

func (nd *BoundLiteralExpressionNode) Type() BoundNodeType {
    return BT_LiteralExpr
}

func (nd *BoundLiteralExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundLiteralExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.LiteralType
} 
