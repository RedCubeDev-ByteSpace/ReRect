package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Loop statement
// --------------
type BoundLoopStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Amount BoundExpressionNode
    Body BoundStatementNode
}

func NewBoundLoopStatementNode(src syntaxnodes.SyntaxNode, amount BoundExpressionNode, body BoundStatementNode) *BoundLoopStatementNode {
    return &BoundLoopStatementNode {
        SourceNode: src,
        Amount: amount,
        Body: body,
    }
}

func (nd *BoundLoopStatementNode) Type() BoundNodeType {
    return BT_LoopStmt
}

func (nd *BoundLoopStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
