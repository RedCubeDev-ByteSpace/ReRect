package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Declaration statement
// ---------------------
type BoundDeclarationStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Variable *symbols.LocalSymbol
    Initializer BoundExpressionNode 
    HasInitializer bool
}

func NewBoundDeclarationStatementNode(src syntaxnodes.SyntaxNode, vari *symbols.LocalSymbol, init BoundExpressionNode, hasinit bool) *BoundDeclarationStatementNode {
    return &BoundDeclarationStatementNode {
        SourceNode: src,
        Variable: vari,
        Initializer: init,
        HasInitializer: hasinit,
    }
}

func (nd *BoundDeclarationStatementNode) Type() BoundNodeType {
    return BT_DeclarationStmt
}

func (nd *BoundDeclarationStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
