package boundnodes

import (
	"bytespace.network/rerect/syntaxnodes"
)

// For statement
// -------------
type BoundForStatementNode struct {
    BoundLoopingStatementNode

    SourceNode syntaxnodes.SyntaxNode

    Initializer BoundStatementNode
    Condition BoundExpressionNode
    Action BoundStatementNode
    Body BoundStatementNode

    BreakLbl BoundLabel
    ContinueLbl BoundLabel
}

func NewBoundForStatementNode(src syntaxnodes.SyntaxNode, init BoundStatementNode, cond BoundExpressionNode, act BoundStatementNode, body BoundStatementNode, brk BoundLabel, cnt BoundLabel) *BoundForStatementNode {
    return &BoundForStatementNode {
        SourceNode: src,
        Initializer: init,
        Condition: cond,
        Action: act,
        Body: body,
        BreakLbl: brk,
        ContinueLbl: cnt,
    }
}

func (nd *BoundForStatementNode) Type() BoundNodeType {
    return BT_ForStmt
}

func (nd *BoundForStatementNode) Source() syntaxnodes.SyntaxNode {
    return nd.SourceNode
}

func (nd *BoundForStatementNode) BreakLabel() BoundLabel {
    return nd.BreakLbl
}

func (nd *BoundForStatementNode) ContinueLabel() BoundLabel {
    return nd.ContinueLbl
}
