package syntaxnodes

import "bytespace.network/rerect/lexer"

// Binary operator precendence
// ---------------------------
func GetBinaryOperatorPrecedence(tok lexer.TokenType) int {
    switch tok {
        case lexer.TT_Star,
             lexer.TT_Slash:
            return 4

        case lexer.TT_Plus,
             lexer.TT_Minus:
            return 3

        case lexer.TT_Ampersands:
            return 2

        case lexer.TT_Pipes:
            return 1

        default:
            return 0
    }
}

// Unary operator precendence 
// ---------------------------
func GetUnaryOperatorPrecedence(tok lexer.TokenType) int {
    switch tok {
        case lexer.TT_Plus,
             lexer.TT_Minus,
             lexer.TT_Bang:
            return 5 // must always be higher than the highest binary op precedence

        default:
            return 0
    }
}
