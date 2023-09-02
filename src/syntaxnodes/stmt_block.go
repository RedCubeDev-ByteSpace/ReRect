package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type BlockStatementNode struct {
    StatementNode

    OpenBrace lexer.Token
    Statements []StatementNode
    CloseBrace lexer.Token
}

func NewBlockStatementNode(openbrace lexer.Token, stmts []StatementNode, closebrace lexer.Token) BlockStatementNode {
    return BlockStatementNode{
        OpenBrace: openbrace,
        Statements: stmts,
        CloseBrace: closebrace,
    }
}

func (n *BlockStatementNode) Position() span.Span {
    return n.OpenBrace.Position.SpanBetween(n.CloseBrace.Position)
}

func (n *BlockStatementNode) Type() SyntaxNodeType {
    return NT_BlockStmt
}
