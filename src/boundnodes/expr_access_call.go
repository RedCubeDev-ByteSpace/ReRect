package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// AccessCall expression
// ---------------------
type BoundAccessCallExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Expression BoundExpressionNode
    Function *symbols.FunctionSymbol
    Arguments []BoundExpressionNode
}

func NewBoundAccessCallExpressionNode(src syntaxnodes.SyntaxNode, exp BoundExpressionNode, fnc *symbols.FunctionSymbol, args []BoundExpressionNode) *BoundAccessCallExpressionNode {
    return &BoundAccessCallExpressionNode {
        SourceNode: src,
        Expression: exp,
        Function: fnc,
        Arguments: args,
    }
}

func (nd *BoundAccessCallExpressionNode) Type() BoundNodeType {
    return BT_AccessCallExpr
}

func (nd *BoundAccessCallExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundAccessCallExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Function.ReturnType
} 
