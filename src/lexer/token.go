// Lexer - token.go
// ---------------------------------------------------------------------
// Contains the defintion for the Token struct aswell as all token types
// ---------------------------------------------------------------------
package lexer

import (
    "bytespace.network/rerect/span"
)

type Token struct {
    Type TokenType
    Buffer string
    Position span.Span
}

func newToken(typ TokenType, buf string, pos span.Span) Token {
    return Token {
        Type: typ,
        Buffer: buf,
        Position: pos,
    }
}

// Types of Tokens (some might even say... token types)
// ----------------------------------------------------
type TokenType string
const (
    // Control Tokens
    TT_WhiteSpace              TokenType = "TT_WhiteSpace"
    TT_EOF                     TokenType = "TT_EOF"
    TT_Comment                 TokenType = "TT_Comment"

    // Punctuation
    TT_Semicolon               TokenType = "TT_Semicolon"
    TT_Colon                   TokenType = "TT_Colon"
    TT_Comma                   TokenType = "TT_Comma"
    TT_OpenParenthesis         TokenType = "TT_OpenParenthesis"
    TT_CloseParenthesis        TokenType = "TT_CloseParenthesis"
    TT_OpenBrackets            TokenType = "TT_OpenBrackets"
    TT_CloseBrackets           TokenType = "TT_CloseBrackets"
    TT_OpenBraces              TokenType = "TT_OpenBraces"
    TT_CloseBraces             TokenType = "TT_CloseBraces"
    
    // Operators 
    TT_LeftArrow               TokenType = "TT_LeftArrow"
    TT_RightArrow              TokenType = "TT_RightArrow"
    TT_Package                 TokenType = "TT_Package"
    
    // Math operators
    TT_Plus                    TokenType = "TT_Plus"
    TT_Minus                   TokenType = "TT_Minus"
    TT_Slash                   TokenType = "TT_Slash"
    TT_Star                    TokenType = "TT_Star"
    TT_Bang                    TokenType = "TT_Bang"
    TT_Equal                   TokenType = "TT_Equal"
    TT_Unequal                 TokenType = "TT_Unequal"
    TT_LessThan                TokenType = "TT_LessEqual"
    TT_GreaterThan             TokenType = "TT_GreaterThan"
    TT_LessEqual               TokenType = "TT_LessEqual"
    TT_GreaterEqual            TokenType = "TT_GreaterEqual"
    TT_Ampersands              TokenType = "TT_Ampersands"
    TT_Pipes                   TokenType = "TT_Pipes"

    // Literals
    TT_String                  TokenType = "TT_String" 
    TT_Integer                 TokenType = "TT_Integer" 
    TT_Float                   TokenType = "TT_Float" 

    // Keywords
    TT_KW_Load                 TokenType = "TT_KW_Load"
    TT_KW_Include              TokenType = "TT_KW_Include"
    TT_KW_Package              TokenType = "TT_KW_Package"
    TT_KW_Function             TokenType = "TT_KW_Function"
    TT_KW_Var                  TokenType = "TT_KW_Var"
    TT_KW_Return               TokenType = "TT_KW_Const"
    TT_KW_While                TokenType = "TT_KW_While"
    TT_KW_From                 TokenType = "TT_KW_From"
    TT_KW_To                   TokenType = "TT_KW_To"
    TT_KW_For                  TokenType = "TT_KW_For"
    TT_KW_Loop                 TokenType = "TT_KW_Loop"
    TT_KW_If                   TokenType = "TT_KW_If"
    TT_KW_Else                 TokenType = "TT_KW_Else"
    TT_KW_True                 TokenType = "TT_KW_True"
    TT_KW_False                TokenType = "TT_KW_False"

    // Identifiers
    TT_Identifier              TokenType = "TT_Identifier"
)

var Keywords = map[string]TokenType {
    "load":     TT_KW_Load,
    "include":  TT_KW_Include,
    "package":  TT_KW_Package,
    "function": TT_KW_Function,
    "var":      TT_KW_Var,
    "return":   TT_KW_Return,
    "while":    TT_KW_While,
    "from":     TT_KW_From,
    "to":       TT_KW_To,
    "for":      TT_KW_For,
    "loop":     TT_KW_Loop,
    "true":     TT_KW_True,
    "false":    TT_KW_False,
    "if":       TT_KW_If,
    "else":     TT_KW_Else,
}

var Symbols = map[string]TokenType {
    "+" : TT_Plus,
    "-" : TT_Minus,
    "/" : TT_Slash,
    "*" : TT_Star,
    "!" : TT_Bang,
    "=" : TT_Equal,
    "!=": TT_Unequal,
    "<" : TT_LessThan,
    ">" : TT_GreaterThan,
    "<=": TT_LessEqual,
    ">=": TT_GreaterEqual,
    "<-": TT_LeftArrow,
    "->": TT_RightArrow,
    "::": TT_Package,
    "&&": TT_Ampersands,
    "||": TT_Pipes,

    "(" : TT_OpenParenthesis,
    ")" : TT_CloseParenthesis,
    "{" : TT_OpenBraces,
    "}" : TT_CloseBraces,
    "[" : TT_OpenBrackets,
    "]" : TT_CloseBrackets,

    "," : TT_Comma,
    ";" : TT_Semicolon,
    ":" : TT_Colon,
}
