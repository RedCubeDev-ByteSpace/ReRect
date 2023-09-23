package syntaxnodes

import "bytespace.network/rerect/span"

// Base node type
// --------------
type SyntaxNode interface {
    Position() span.Span
    Type() SyntaxNodeType 
}

// Member node type
// ----------------
type MemberNode interface {
   SyntaxNode 
}

// Statement node type
// -------------------
type StatementNode interface {
    SyntaxNode
}

// Expression node type
// --------------------
type ExpressionNode interface {
    SyntaxNode
}

// Node types
// ----------
type SyntaxNodeType string
const (
    // Members
    NT_Load               SyntaxNodeType = "Load member"
    NT_Package            SyntaxNodeType = "Package member"
    NT_Function           SyntaxNodeType = "Function member"
    NT_Global             SyntaxNodeType = "Global variable member"
    NT_Container          SyntaxNodeType = "Container member"
    NT_Trait              SyntaxNodeType = "Trait member"

    // Statements
    NT_DeclarationStmt    SyntaxNodeType = "Local variable statement"
    NT_ReturnStmt         SyntaxNodeType = "Return statement"
    NT_WhileStmt          SyntaxNodeType = "While statement"
    NT_FromToStmt         SyntaxNodeType = "From-To statement node"
    NT_ForStmt            SyntaxNodeType = "For statement node"
    NT_LoopStmt           SyntaxNodeType = "Loop statement node"
    NT_BreakStmt          SyntaxNodeType = "Break statement node"
    NT_ContinueStmt       SyntaxNodeType = "Continue statement node"
    NT_BlockStmt          SyntaxNodeType = "Block statement node"
    NT_ExpressionStmt     SyntaxNodeType = "Expression statement node"
    NT_IfStmt             SyntaxNodeType = "If statement node"

    // Expressions
    NT_LiteralExpr        SyntaxNodeType = "Literal expression node"
    NT_AssignmentExpr     SyntaxNodeType = "Assignment expression node"
    NT_UnaryExpr          SyntaxNodeType = "Unary expression node"
    NT_BinaryExpr         SyntaxNodeType = "Binary expression node"
    NT_CallExpr           SyntaxNodeType = "Call expression node"
    NT_NameExpr           SyntaxNodeType = "Name expression node"
    NT_ParenthesizedExpr  SyntaxNodeType = "Parenthesized expression node"
    NT_MakeArrayExpr      SyntaxNodeType = "Array creation expression node"
    NT_ArrayIndexExpr     SyntaxNodeType = "Array index expression node"
    NT_AccessExpr         SyntaxNodeType = "Access expression node"
    NT_MakeExpr           SyntaxNodeType = "Object creation expression node"

    NT_ErrorExpr          SyntaxNodeType = "Error expression node"

    // Clauses
    NT_ParameterCls       SyntaxNodeType = "Parameter clause"
    NT_TypeCls            SyntaxNodeType = "Type clause"
    NT_FieldCls           SyntaxNodeType = "Field clause"
    NT_FieldAssignmentCls SyntaxNodeType = "Field assignment clause"
    NT_TraitCls           SyntaxNodeType = "Trait clause"
)
