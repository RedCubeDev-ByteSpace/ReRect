package boundnodes

import (
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/symbols"
)

// Binary operators
//----------------
type BoundBinaryOperator struct {
    Left  *symbols.TypeSymbol
    Right *symbols.TypeSymbol
    Operation BinaryOperatorType
    Result *symbols.TypeSymbol
}

func NewBoundBinaryOperator(op BinaryOperatorType, left *symbols.TypeSymbol, right *symbols.TypeSymbol, result *symbols.TypeSymbol) *BoundBinaryOperator {
    return &BoundBinaryOperator{
        Left: left,
        Right: right,
        Operation: op,
        Result: result,
    }
}

// List of unary operations
// ------------------------
type BinaryOperatorType string;
const (
    BO_Addition       BinaryOperatorType = "Addition operator"
    BO_Subtraction    BinaryOperatorType = "Subtraction operator"
    BO_Multiplication BinaryOperatorType = "Multiplication operator"
    BO_Division       BinaryOperatorType = "Division operator"

    BO_LogicalAnd     BinaryOperatorType = "Logical and operator"
    BO_LogicalOr      BinaryOperatorType = "Logical or operator"

    BO_Equal          BinaryOperatorType = "Equality operator"
    BO_UnEqual        BinaryOperatorType = "Unequality operator"
    BO_LessThan       BinaryOperatorType = "Less than operator"
    BO_LessEqual      BinaryOperatorType = "Less equal operator"
    BO_GreaterThan    BinaryOperatorType = "Greater than operator"
    BO_GreaterEqual   BinaryOperatorType = "Greater equal operator"

    BO_Concat         BinaryOperatorType = "String concat operator"
)

func GetBinaryOperator(op lexer.TokenType, left *symbols.TypeSymbol, right *symbols.TypeSymbol) *BoundBinaryOperator {
    // Basic arithmetic (plus, minus, multiply, divide)
    // Comparing (=, !=, <, <=, >, >=)
    if (op == lexer.TT_Plus          ||
        op == lexer.TT_Minus         ||
        op == lexer.TT_Star          ||
        op == lexer.TT_Slash         ||
        op == lexer.TT_Equal         ||
        op == lexer.TT_Unequal       ||
        op == lexer.TT_LessThan      ||
        op == lexer.TT_LessEqual     ||
        op == lexer.TT_GreaterThan   ||
        op == lexer.TT_GreaterEqual) && (
        (left.TypeGroup == symbols.INT   && right.TypeGroup == symbols.INT) ||
        (left.TypeGroup == symbols.FLOAT && right.TypeGroup == symbols.FLOAT)){ 
    
        // always use the larger type
        typ := left
        if (right.TypeSize > typ.TypeSize) {
            typ = right
        }
       
        // mmmm operations
        switch op {
        case lexer.TT_Plus: 
            return NewBoundBinaryOperator(BO_Addition, typ, typ, typ)
        case lexer.TT_Minus: 
            return NewBoundBinaryOperator(BO_Subtraction, typ, typ, typ)
        case lexer.TT_Star: 
            return NewBoundBinaryOperator(BO_Multiplication, typ, typ, typ)
        case lexer.TT_Slash: 
            return NewBoundBinaryOperator(BO_Division, typ, typ, typ)
        case lexer.TT_Equal: 
            return NewBoundBinaryOperator(BO_Equal, typ, typ, typ)
        case lexer.TT_Unequal: 
            return NewBoundBinaryOperator(BO_UnEqual, typ, typ, typ)
        case lexer.TT_LessThan: 
            return NewBoundBinaryOperator(BO_LessThan, typ, typ, typ)
        case lexer.TT_LessEqual: 
            return NewBoundBinaryOperator(BO_LessEqual, typ, typ, typ)
        case lexer.TT_GreaterThan: 
            return NewBoundBinaryOperator(BO_GreaterThan, typ, typ, typ)
        case lexer.TT_GreaterEqual: 
            return NewBoundBinaryOperator(BO_GreaterEqual, typ, typ, typ)
        }
    }

    // Logical operations (and, or)
    if (op == lexer.TT_Ampersands ||
        op == lexer.TT_Pipes    ) && (
            left.Equal(compunit.GlobalDataTypeRegister["bool"]) &&
            right.Equal(compunit.GlobalDataTypeRegister["bool"])){ 
      
        typ := compunit.GlobalDataTypeRegister["bool"]

        switch op {
        case lexer.TT_Ampersands: 
            return NewBoundBinaryOperator(BO_LogicalAnd, typ, typ, typ)
        case lexer.TT_Pipes: 
            return NewBoundBinaryOperator(BO_LogicalOr, typ, typ, typ)
        }
    }

    // Equality and unequality
    if left.Equal(right) && op == lexer.TT_Equal {
        return NewBoundBinaryOperator(BO_Equal, left, right, compunit.GlobalDataTypeRegister["bool"])
    } 

    if left.Equal(right) && op == lexer.TT_Unequal {
        return NewBoundBinaryOperator(BO_UnEqual, left, right, compunit.GlobalDataTypeRegister["bool"])
    } 

    // string concat
    if left.Equal(compunit.GlobalDataTypeRegister["string"]) &&
       right.Equal(compunit.GlobalDataTypeRegister["string"]) {
        return NewBoundBinaryOperator(BO_Concat, left, right, left) 
    }

    return nil
}

