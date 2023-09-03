package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Expression statement
// --------------------
type BoundExpressionStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Expression BoundExpressionNode
}

func NewBoundExpressionStatementNode(src syntaxnodes.SyntaxNode, expr BoundExpressionNode) *BoundExpressionStatementNode {
    return &BoundExpressionStatementNode {
        SourceNode: src,
        Expression: expr,
    }
}

func (nd *BoundExpressionStatementNode) Type() BoundNodeType {
    return BT_ExpressionStmt
}

func (nd *BoundExpressionStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
