// Lexer - lexer.go
// ---------------------------------------------------------------------
// Contains most things lexing
// ---------------------------------------------------------------------
package lexer

import (
	"os"
	"strings"
	"unicode"

	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/span"
)

// Lexer struct
// ------------
type Lexer struct {
    SourceStr    string
    Source       []rune
    SourceFileId int
   
    Index  int
    Length int

    Tokens []Token
}

const EOF rune = '\004'

// --------------------------------------------------------
// Helpers
// --------------------------------------------------------
func (lxr *Lexer) current() rune {
    return lxr.peek(0)
}

func (lxr *Lexer) peek(offset int) rune {
    if offset + lxr.Index >= lxr.Length {
        return EOF
    }

    return lxr.Source[offset + lxr.Index]
}

func (lxr *Lexer) append(tok Token) {
    if tok.Type == TT_WhiteSpace || tok.Type == TT_Comment {
        return
    }

    lxr.Tokens = append(lxr.Tokens, tok)
} 

func (lxr *Lexer) currentSpan() span.Span {
    return span.Span{
        File: lxr.SourceFileId,
        FromIdx: lxr.Index,
        ToIdx: lxr.Index,
    }
}

func (lxr *Lexer) step(size int) {
    lxr.Index += size
}

// --------------------------------------------------------
// Lexing
// --------------------------------------------------------

// Lex a given file
// ----------------
func LexFile(file string) []Token {
    txt, err := os.ReadFile(file)

    // Report IO Errors
    if err != nil {
         error.Report(
            error.NewError(error.FIO, span.Internal(), err.Error()),
         )

         return make([]Token, 0)
    }

    // Remember this file
    idx := compunit.RegisterSource(file, string(txt))
    
    return LexString(string(txt), idx)
}

// Lex a given string
// ------------------
func LexString(code string, srcidx int) []Token {
   // Instantiate a new lexer
   // -----------------------
   lex := Lexer {
       Source: []rune(code),
       SourceStr: code,
       SourceFileId: srcidx,

       Length: len(code),

       Tokens: make([]Token, 0),
   }

   lex.lex();

   return lex.Tokens
}

func (lxr *Lexer) lex() {
    for true {
        current := lxr.current()
 
        // We have reached the end of file
        if current == EOF {
            break
        }
       
        // Space, Tab, Newline, Carriage Return -> Whitespace
        if current == ' ' || current == '\t' || current == '\n' || current == '\r' {
            lxr.append(Token{
                Type: TT_WhiteSpace,
                Buffer: string(current),
                Position: lxr.currentSpan(),
            })

            lxr.step(1)

        // Double slash -> Comment
        } else if current == '/' && lxr.peek(1) == '/' {
            lxr.lexComment()

        // Quotes -> String
        } else if current == '"' {
            lxr.lexString()

        // Digit -> Number
        } else if unicode.IsDigit(current) {
            lxr.lexNumber()

        // Letter -> Keyword or Identifier
        } else if unicode.IsLetter(current) {
            lxr.lexWord()
        
        // Otherwise -> probably some symbol or operator
        } else {
            lxr.lexSymbol()
        }
    }

    // Append EOF token at the end of the token list
    lxr.append(Token{
        Type: TT_EOF,
        Position: lxr.currentSpan(),
    })
}

// String lexing
// -------------
func (lxr *Lexer) lexString() {
    startPos := lxr.currentSpan()
    
    // step over the leading qoute
    lxr.step(1)

    // the string content buffer
    buffer := ""

    // As long as we dont find the EOF or a closing qoute
    for lxr.current() != EOF && lxr.current() != '"' {
        buffer += string(lxr.current())
        lxr.step(1)
    }

    // assemble the token
    tok := Token {
        Type: TT_String,
        Buffer: buffer,
        Position: startPos.SpanBetween(lxr.currentSpan()),
    }

    // store the token
    lxr.append(tok)
    
    // step over the trailing qoute
    lxr.step(1)
}

// Comment lexing
// --------------
func (lxr *Lexer) lexComment() {
    startPos := lxr.currentSpan()

    // Step forward until we hit an EOF or an end of line
    for lxr.current() != EOF && lxr.current() != '\n' {
        lxr.step(1)
    }


    // assemble the token
    tok := Token {
        Type: TT_Comment,
        Position: startPos.SpanBetween(lxr.currentSpan()),
    }

    // store the token
    lxr.append(tok)
}

// Number lexing
// -------------
func (lxr *Lexer) lexNumber() {
    // remember the start of this number
    startPos := lxr.currentSpan()
    
    // the number buffer
    buffer := ""

    // is this a floating point number?
    hasDecimal := false

    // Fill buffer as long as we conitnue to read digits or decimal points
    for unicode.IsDigit(lxr.current()) || lxr.current() == '.' || lxr.current() == '_' {
        // Skip underscores (they are just for visual clarity)
        if lxr.current() == '_' {
            lxr.step(1)
            continue
        }

        // If we found a decimal point -> this is a floating point number
        if lxr.current() == '.' {
            
            // We already found a decimal point (multiple decimal points are illegal)
            if hasDecimal {
                error.Report(error.NewError(error.LEX, lxr.currentSpan(), "Illegal decimal point!"))
                return
            }

            // otherwise...
            hasDecimal = true
        }

        // add digit to buffer
        buffer += string(lxr.current())
        lxr.step(1)
    }

    var tok Token

    if hasDecimal {
        // assemble the token
        tok = Token {
            Type: TT_Float,
            Buffer: buffer,
            Position: startPos.SpanBetween(lxr.currentSpan()),
        }
    } else {
        // assemble the token
        tok = Token {
            Type: TT_Integer,
            Buffer: buffer,
            Position: startPos.SpanBetween(lxr.currentSpan()),
        }
    }
    
    // store the token
    lxr.append(tok)
}

// Keywords and identifierts
// -------------------------
func (lxr *Lexer) lexWord() {
    // remember current position
    startPos := lxr.currentSpan()

    // keyword buffer
    buffer := ""

    for unicode.IsLetter(lxr.current()) || unicode.IsDigit(lxr.current()) || lxr.current() == '_' {
        buffer += string(lxr.current())
        lxr.step(1)
    } 

    tokenType := lxr.identifyKeyword(buffer)

    // assemble the token
    tok := Token {
        Type: tokenType,
        Buffer: buffer,
        Position: startPos.SpanBetween(lxr.currentSpan()),
    }

    // store the token
    lxr.append(tok)
}

// Find token type for keywords
// ----------------------------
func (lxr *Lexer) identifyKeyword(buf string) TokenType {
    typ, ok := Keywords[buf]

    // we found nothing -> this is an identifier
    if !ok {
        return TT_Identifier

    // we found something -> we found something (lmao)
    } else {
        return typ
    }
}

// Lex operators and other symbols
// -------------------------------
func (lxr *Lexer) lexSymbol() {
    // remember current position
    startPos := lxr.currentSpan()

    // operator buffer
    buffer := ""

    // which token could this be?
    var possibleToken TokenType = TT_EOF
    tokText := ""

    for true {
        buffer += string(lxr.current())

        found := false
        for k, v := range Symbols {

            // we found at least one operator that starts with the current buffer
            if strings.HasPrefix(k, buffer) {

                // did we already find a match?
                if found {
                    // if this token is shorter than the last -> use this one
                    if len(k) < len(tokText) {
                        // store the match
                        possibleToken = v
                        tokText = k
                    }

                // if we havent found anything yet -> use this
                } else {
                    possibleToken = v
                    tokText = k
                }

                // flag to continue to the next char
                found = true
            }
        }

        // continue to next char
        if found {
            // step to the next char
            lxr.step(1)

            continue
        }

        // if we found nothing -> did we find one before?
        if possibleToken != TT_EOF {
            // return that one
            lxr.append(Token {
                Type: possibleToken,
                Position: startPos.SpanBetween(lxr.currentSpan()),
            })

            return
        }

        // if not -> theres no token that matches the input
        error.Report(error.NewError(error.LEX, startPos.SpanBetween(lxr.currentSpan()), "Unknown symbol!"))

        // step over the char
        lxr.step(1)

        return
    }
}
