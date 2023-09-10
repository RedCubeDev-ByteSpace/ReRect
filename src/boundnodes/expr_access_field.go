package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// AccessField expression
// ---------------------
type BoundAccessFieldExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Expression BoundExpressionNode
    Field *symbols.FieldSymbol
}

func NewBoundAccessFieldExpressionNode(src syntaxnodes.SyntaxNode, exp BoundExpressionNode, fld *symbols.FieldSymbol) *BoundAccessFieldExpressionNode {
    return &BoundAccessFieldExpressionNode {
        SourceNode: src,
        Expression: exp,
        Field: fld,
    }
}

func (nd *BoundAccessFieldExpressionNode) Type() BoundNodeType {
    return BT_AccessFieldExpr
}

func (nd *BoundAccessFieldExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundAccessFieldExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Field.VarType()
} 
