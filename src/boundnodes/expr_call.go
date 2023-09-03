package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Call expression
// ---------------
type BoundCallExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Function *symbols.FunctionSymbol
    Arguments []BoundExpressionNode
}

func NewBoundCallExpressionNode(src syntaxnodes.SyntaxNode, fnc *symbols.FunctionSymbol, args []BoundExpressionNode) *BoundCallExpressionNode {
    return &BoundCallExpressionNode {
        SourceNode: src,
        Function: fnc,
        Arguments: args,
    }
}

func (nd *BoundCallExpressionNode) Type() BoundNodeType {
    return BT_CallExpr
}

func (nd *BoundCallExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundCallExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Function.ReturnType
} 
