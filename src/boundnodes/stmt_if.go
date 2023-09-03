package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// If statement
// ------------
type BoundIfStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Condition BoundExpressionNode
    Body BoundStatementNode

    ElseBody BoundStatementNode
    HasElse bool
}

func NewBoundIfStatementNode(src syntaxnodes.SyntaxNode, cond BoundExpressionNode, body BoundStatementNode, elsebody BoundStatementNode, haselse bool) *BoundIfStatementNode {
    return &BoundIfStatementNode {
        SourceNode: src,
        Condition: cond,
        Body: body,
        ElseBody: elsebody,
        HasElse: haselse,
    }
}

func (nd *BoundIfStatementNode) Type() BoundNodeType {
    return BT_IfStmt
}

func (nd *BoundIfStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
