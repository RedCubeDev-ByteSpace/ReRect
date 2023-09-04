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

    BreakLbl BoundLabel
    ContinueLbl BoundLabel
}

func NewBoundFromToStatementNode(src syntaxnodes.SyntaxNode, iterator symbols.VariableSymbol, lb BoundExpressionNode, up BoundExpressionNode, body BoundStatementNode, brk BoundLabel, cnt BoundLabel) *BoundFromToStatementNode {
    return &BoundFromToStatementNode {
        SourceNode: src,
        Iterator: iterator,
        LowerBound: lb,
        UpperBound: up,
        Body: body,
        BreakLbl: brk,
        ContinueLbl: cnt,
    }
}

func (nd *BoundFromToStatementNode) Type() BoundNodeType {
    return BT_FromToStmt
}

func (nd *BoundFromToStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundFromToStatementNode) BreakLabel() BoundLabel {
    return nd.BreakLbl
}

func (nd *BoundFromToStatementNode) ContinueLabel() BoundLabel {
    return nd.ContinueLbl
}
