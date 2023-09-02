package syntaxnodes

import (
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/span"
)

type PackageNode struct {
    MemberNode

    PackageKw lexer.Token
    PackageName lexer.Token
}

func NewPackageNode(packagekw lexer.Token, pck lexer.Token) *PackageNode {
    return &PackageNode{
        PackageKw: packagekw,
        PackageName: pck,
    }
}

func (n *PackageNode) Position() span.Span {
    return n.PackageKw.Position.SpanBetween(n.PackageName.Position)
}

func (n *PackageNode) Type() SyntaxNodeType {
    return NT_Package
}
