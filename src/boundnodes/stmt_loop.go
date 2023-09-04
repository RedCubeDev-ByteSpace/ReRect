package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// Loop statement
// --------------
type BoundLoopStatementNode struct {
    BoundLoopingStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Amount BoundExpressionNode
    Body BoundStatementNode

    BreakLbl BoundLabel
    ContinueLbl BoundLabel
}

func NewBoundLoopStatementNode(src syntaxnodes.SyntaxNode, amount BoundExpressionNode, body BoundStatementNode, brk BoundLabel, cnt BoundLabel) *BoundLoopStatementNode {
    return &BoundLoopStatementNode {
        SourceNode: src,
        Amount: amount,
        Body: body,
        BreakLbl: brk,
        ContinueLbl: cnt,
    }
}

func (nd *BoundLoopStatementNode) Type() BoundNodeType {
    return BT_LoopStmt
}

func (nd *BoundLoopStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundLoopStatementNode) BreakLabel() BoundLabel {
    return nd.BreakLbl
}

func (nd *BoundLoopStatementNode) ContinueLabel() BoundLabel {
    return nd.ContinueLbl
}
