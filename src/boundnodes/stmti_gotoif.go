package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Goto statement
// --------------
type BoundGotoIfStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Label string
    Condition BoundExpressionNode
}

func NewBoundGotoIfStatementNode(src syntaxnodes.SyntaxNode, label string, cond BoundExpressionNode) *BoundGotoIfStatementNode {
    return &BoundGotoIfStatementNode {
        SourceNode: src,
        Label: label,
        Condition: cond,
    }
}

func (nd *BoundGotoIfStatementNode) Type() BoundNodeType {
    return BT_GoToIfIStmt
}

func (nd *BoundGotoIfStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
