package boundnodes

import (
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/symbols"
)

// Unary operators
//----------------
type BoundUnaryOperator struct {
    Operand *symbols.TypeSymbol
    Operation UnaryOperatorType
    Result *symbols.TypeSymbol
}

func NewBoundUnaryOperator(op UnaryOperatorType, operand *symbols.TypeSymbol, result *symbols.TypeSymbol) *BoundUnaryOperator {
    return &BoundUnaryOperator{
        Operand: operand,
        Operation: op,
        Result: result,
    }
}

// List of unary operations
// ------------------------
type UnaryOperatorType string;
const (
    UO_Identity        UnaryOperatorType = "Identity operator"
    UO_Negation        UnaryOperatorType = "Negation operator"
    UO_LogicalNegation UnaryOperatorType = "Logical negation operator"
)

func GetUnaryOperator(op lexer.TokenType, operand *symbols.TypeSymbol) *BoundUnaryOperator {
    if op == lexer.TT_Plus && (
        operand.TypeGroup == symbols.INT ||
        operand.TypeGroup == symbols.FLOAT) {
        
        return NewBoundUnaryOperator(UO_Identity, operand, operand)
    }

    if op == lexer.TT_Minus && (
        operand.TypeGroup == symbols.INT ||
        operand.TypeGroup == symbols.FLOAT) {
        
        return NewBoundUnaryOperator(UO_Negation, operand, operand)
    }

    if op == lexer.TT_Bang && operand.Equal(compunit.GlobalDataTypeRegister["bool"]) {
        return NewBoundUnaryOperator(UO_LogicalNegation, operand, operand)
    }

    return nil
}

