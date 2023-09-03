package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Assignment expression
// ---------------------
type BoundAssignmentExpressionNode struct {
    BoundExpressionNode

    SourceNode syntaxnodes.SyntaxNode

    Variable symbols.VariableSymbol
    Value BoundExpressionNode
}

func NewBoundAssignmentExpressionNode(src syntaxnodes.SyntaxNode, vari symbols.VariableSymbol, val BoundExpressionNode) *BoundAssignmentExpressionNode {
    return &BoundAssignmentExpressionNode {
        SourceNode: src,
        Variable: vari,
        Value: val,
    }
}

func (nd *BoundAssignmentExpressionNode) Type() BoundNodeType {
    return BT_AssignmentExpr
}

func (nd *BoundAssignmentExpressionNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundAssignmentExpressionNode) ExprType() *symbols.TypeSymbol {
    return nd.Variable.VarType()
} 
