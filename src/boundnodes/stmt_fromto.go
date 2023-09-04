package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// FromTo statement
// ----------------
type BoundFromToStatementNode struct {
    BoundStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Iterator symbols.VariableSymbol
    LowerBound BoundExpressionNode
    UpperBound BoundExpressionNode
    Body BoundStatementNode
}

func NewBoundFromToStatementNode(src syntaxnodes.SyntaxNode, iterator symbols.VariableSymbol, lb BoundExpressionNode, up BoundExpressionNode, body BoundStatementNode) *BoundFromToStatementNode {
    return &BoundFromToStatementNode {
        SourceNode: src,
        Iterator: iterator,
        LowerBound: lb,
        UpperBound: up,
        Body: body,
    }
}

func (nd *BoundFromToStatementNode) Type() BoundNodeType {
    return BT_FromToStmt
}

func (nd *BoundFromToStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}
