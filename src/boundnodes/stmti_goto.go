package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Goto statement
// --------------
type BoundGotoStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Label string
}

func NewBoundGotoStatementNode(src syntaxnodes.SyntaxNode, label string) *BoundGotoStatementNode {
    return &BoundGotoStatementNode {
        SourceNode: src,
        Label: label,
    }
}

func (nd *BoundGotoStatementNode) Type() BoundNodeType {
    return BT_GoToIStmt
}

func (nd *BoundGotoStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
