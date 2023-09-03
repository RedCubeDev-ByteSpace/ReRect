package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Block statement
// --------------
type BoundBlockStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Statements []BoundStatementNode
}

func NewBoundBlockStatementNode(src syntaxnodes.SyntaxNode, stmts []BoundStatementNode) *BoundBlockStatementNode {
    return &BoundBlockStatementNode {
        SourceNode: src,
        Statements: stmts,
    }
}

func (nd *BoundBlockStatementNode) Type() BoundNodeType {
    return BT_BlockStmt
}

func (nd *BoundBlockStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
