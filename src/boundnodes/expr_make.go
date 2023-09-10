package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Make expression
// ------------------
type BoundMakeExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Container *symbols.ContainerSymbol
    
    Initializer map[*symbols.FieldSymbol]BoundExpressionNode
    HasInitializer bool

    Arguments []BoundExpressionNode
    HasConstructor bool
}

func NewBoundMakeExpressionNode(src syntaxnodes.SyntaxNode, cnt *symbols.ContainerSymbol, init map[*symbols.FieldSymbol]BoundExpressionNode, hasinit bool, args []BoundExpressionNode, hascst bool) *BoundMakeExpressionNode {
    return &BoundMakeExpressionNode {
        SourceNode: src,
        Container: cnt,

        Initializer: init,
        HasInitializer: hasinit,

        Arguments: args,
        HasConstructor: hascst,
    }
}

func (nd *BoundMakeExpressionNode) Type() BoundNodeType {
    return BT_MakeExpr
}

func (nd *BoundMakeExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundMakeExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Container.ContainerType
} 
