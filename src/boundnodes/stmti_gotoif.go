package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Goto statement
// --------------
type BoundGotoIfStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Label BoundLabel
    Condition BoundExpressionNode
}

func NewBoundGotoIfStatementNode(src syntaxnodes.SyntaxNode, label BoundLabel, cond BoundExpressionNode) *BoundGotoIfStatementNode {
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
