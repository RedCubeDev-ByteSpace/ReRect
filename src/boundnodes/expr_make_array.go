package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// MakeArray expression
// ------------------
type BoundMakeArrayExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    ArrType *symbols.TypeSymbol
    Length BoundExpressionNode
    Initializer []BoundExpressionNode
    HasInitializer bool
}

func NewBoundMakeArrayExpressionNode(src syntaxnodes.SyntaxNode, arrtyp *symbols.TypeSymbol, length BoundExpressionNode, init []BoundExpressionNode, hasinit bool) *BoundMakeArrayExpressionNode {
    return &BoundMakeArrayExpressionNode {
        SourceNode: src,
        ArrType: arrtyp,
        Length: length,
        Initializer: init,
        HasInitializer: hasinit,
    }
}

func (nd *BoundMakeArrayExpressionNode) Type() BoundNodeType {
    return BT_MakeArrayExpr
}

func (nd *BoundMakeArrayExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundMakeArrayExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.ArrType
} 
