package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// While statement
// ---------------
type BoundWhileStatementNode struct {
    BoundLoopingStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Condtion BoundExpressionNode
    Body BoundStatementNode

    BreakLbl BoundLabel
    ContinueLbl BoundLabel
}

func NewBoundWhileStatementNode(src syntaxnodes.SyntaxNode, cond BoundExpressionNode, body BoundStatementNode, brk BoundLabel, cnt BoundLabel) *BoundWhileStatementNode {
    return &BoundWhileStatementNode {
        SourceNode: src,
        Condtion: cond,
        Body: body,
        BreakLbl: brk,
        ContinueLbl: cnt,
    }
}

func (nd *BoundWhileStatementNode) Type() BoundNodeType {
    return BT_WhileStmt
}

func (nd *BoundWhileStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundWhileStatementNode) BreakLabel() BoundLabel {
    return nd.BreakLbl
}

func (nd *BoundWhileStatementNode) ContinueLabel() BoundLabel {
    return nd.ContinueLbl
}
