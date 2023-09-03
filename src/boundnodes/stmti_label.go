package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Label statement
// ---------------
type BoundLabelStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Label string
}

func NewBoundLabelStatementNode(src syntaxnodes.SyntaxNode, label string) *BoundLabelStatementNode {
    return &BoundLabelStatementNode {
        SourceNode: src,
        Label: label,
    }
}

func (nd *BoundLabelStatementNode) Type() BoundNodeType {
    return BT_LabelIStmt
}

func (nd *BoundLabelStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
