package boundnodes

import (
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Error expression
// ----------------
type BoundErrorExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode
}

func NewBoundErrorExpressionNode(src syntaxnodes.SyntaxNode) *BoundErrorExpressionNode {
    return &BoundErrorExpressionNode {
        SourceNode: src,
    }
}

func (nd *BoundErrorExpressionNode) Type() BoundNodeType {
    return BT_ErrorExpr
}

func (nd *BoundErrorExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundErrorExpressionNode) ExprType() *symbols.TypeSymbol {
    return compunit.GlobalDataTypeRegister["error"]
} 
