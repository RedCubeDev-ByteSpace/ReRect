package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// For statement
// -------------
type BoundForStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Initializer BoundStatementNode
    Condition BoundExpressionNode
    Action BoundStatementNode
    Body BoundStatementNode
}

func NewBoundForStatementNode(src syntaxnodes.SyntaxNode, init BoundStatementNode, cond BoundExpressionNode, act BoundStatementNode, body BoundStatementNode) *BoundForStatementNode {
    return &BoundForStatementNode {
        SourceNode: src,
        Initializer: init,
        Condition: cond,
        Action: act,
        Body: body,
    }
}

func (nd *BoundForStatementNode) Type() BoundNodeType {
    return BT_ForStmt
}

func (nd *BoundForStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
