package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Return statement
// ---------------------
type BoundReturnStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    ReturnValue BoundExpressionNode 
    HasReturnValue bool
}

func NewBoundReturnStatementNode(src syntaxnodes.SyntaxNode, retv BoundExpressionNode, hasretv bool) *BoundReturnStatementNode {
    return &BoundReturnStatementNode {
        SourceNode: src,
        ReturnValue: retv,
        HasReturnValue: hasretv,
    }
}

func (nd *BoundReturnStatementNode) Type() BoundNodeType {
    return BT_ReturnStmt
}

func (nd *BoundReturnStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
