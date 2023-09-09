package boundnodes

import (
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

// Base bound nodes type
// ---------------------
type BoundNode interface {
    Source() syntaxnodes.SyntaxNode
    Type() BoundNodeType
}

// Bound statement nodes
// ---------------------
type BoundStatementNode interface {
    BoundNode
}

// Bound loop statement node
// -------------------------
type BoundLoopingStatementNode interface {
    BoundStatementNode
    BreakLabel()    BoundLabel
    ContinueLabel() BoundLabel
}

type BoundLabel string

// Bound expression nodes
// ----------------------
type BoundExpressionNode interface {
    BoundNode
    ExprType() *symbols.TypeSymbol
}

// Bound node types
// ----------------
type BoundNodeType string;
const (
    // Statements
    BT_DeclarationStmt BoundNodeType = "Local variable declaration"
    BT_ReturnStmt      BoundNodeType = "Return statement"
    BT_WhileStmt       BoundNodeType = "While statement"
    BT_FromToStmt      BoundNodeType = "From-To statement"
    BT_ForStmt         BoundNodeType = "For statement"
    BT_LoopStmt        BoundNodeType = "Loop statement"
    BT_BlockStmt       BoundNodeType = "Block statement"
    BT_ExpressionStmt  BoundNodeType = "Expression statement"
    BT_IfStmt          BoundNodeType = "If statement"

    // Internal VM statements
    BT_LabelIStmt      BoundNodeType = "Internal label statement"
    BT_GoToIStmt       BoundNodeType = "Internal goto statement"
    BT_GoToIfIStmt     BoundNodeType = "Internal conditional goto statement"
    BT_DeleteIStmt     BoundNodeType = "Internal variable deletion statement"
    BT_ApproachIStmt   BoundNodeType = "Internal increment / decrement statement"

    // Expressions
    BT_LiteralExpr     BoundNodeType = "Literal expression"
    BT_AssignmentExpr  BoundNodeType = "Assignment expression"
    BT_UnaryExpr       BoundNodeType = "Unary expression"
    BT_BinaryExpr      BoundNodeType = "Binary expression"
    BT_CallExpr        BoundNodeType = "Call expression"
    BT_NameExpr        BoundNodeType = "Name expression"
    BT_ConversionExpr  BoundNodeType = "Conversion expression"
    BT_MakeArrayExpr   BoundNodeType = "Array creation expression"
    BT_ArrayIndexExpr  BoundNodeType = "Array index expression"

    BT_ErrorExpr       BoundNodeType = "Error expression"
)
