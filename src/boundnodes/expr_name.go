package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Name expression
// ---------------
type BoundNameExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Variable symbols.VariableSymbol
}

func NewBoundNameExpressionNode(src syntaxnodes.SyntaxNode, vari symbols.VariableSymbol) *BoundNameExpressionNode {
    return &BoundNameExpressionNode {
        SourceNode: src,
        Variable: vari,
    }
}

func (nd *BoundNameExpressionNode) Type() BoundNodeType {
    return BT_NameExpr
}

func (nd *BoundNameExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundNameExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Variable.VarType()
} 
