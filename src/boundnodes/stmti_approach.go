package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Approach statement
// ----------------
type BoundApproachStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Iterator symbols.VariableSymbol
    Target BoundExpressionNode
}

func NewBoundApproachStatementNode(src syntaxnodes.SyntaxNode, vari symbols.VariableSymbol, target BoundExpressionNode) *BoundApproachStatementNode {
    return &BoundApproachStatementNode {
        SourceNode: src,
        Iterator: vari,
        Target: target,
    }
}

func (nd *BoundApproachStatementNode) Type() BoundNodeType {
    return BT_ApproachIStmt
}

func (nd *BoundApproachStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
