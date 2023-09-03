package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// While statement
// ---------------
type BoundWhileStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Condtion BoundExpressionNode
    Body BoundStatementNode
}

func NewBoundWhileStatementNode(src syntaxnodes.SyntaxNode, cond BoundExpressionNode, body BoundStatementNode) *BoundWhileStatementNode {
    return &BoundWhileStatementNode {
        SourceNode: src,
        Condtion: cond,
        Body: body,
    }
}

func (nd *BoundWhileStatementNode) Type() BoundNodeType {
    return BT_WhileStmt
}

func (nd *BoundWhileStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
