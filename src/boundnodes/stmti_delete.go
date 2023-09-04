package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Delete statement
// --------------
type BoundDeleteStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Variable symbols.VariableSymbol
}

func NewBoundDeleteStatementNode(src syntaxnodes.SyntaxNode, vari symbols.VariableSymbol) *BoundDeleteStatementNode {
    return &BoundDeleteStatementNode {
        SourceNode: src,
        Variable: vari,
    }
}

func (nd *BoundDeleteStatementNode) Type() BoundNodeType {
    return BT_DeleteIStmt
}

func (nd *BoundDeleteStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
