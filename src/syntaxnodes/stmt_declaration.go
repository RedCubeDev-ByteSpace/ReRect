package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type DeclarationStatementNode struct {
    StatementNode

    VarKw lexer.Token
    VarName lexer.Token

    VarType TypeClauseNode
    HasExplicitType bool

    Initializer ExpressionNode
    HasInitializer bool
}

func NewDeclarationStatementNode(varkw lexer.Token, varname lexer.Token, typ TypeClauseNode, hastyp bool, init ExpressionNode, hasinit bool) DeclarationStatementNode {
    return DeclarationStatementNode{
        VarKw: varkw,
        VarName: varname,
        VarType: typ,
        HasExplicitType: hastyp,
        Initializer: init,
        HasInitializer: hasinit,
    }
}

func (n *DeclarationStatementNode) Position() span.Span {
    spn := n.VarKw.Position

    if n.HasExplicitType {
        spn = spn.SpanBetween(n.VarType.Position())
    }

    if n.HasInitializer {
        spn = spn.SpanBetween(n.Initializer.Position())
    }

    return spn
}

func (n *DeclarationStatementNode) Type() SyntaxNodeType {
    return NT_DeclarationStmt
}
