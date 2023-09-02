// Error - error.go
// ---------------------------------------------------------------------
// Contains utilities for error handling and diagnostics
// ---------------------------------------------------------------------
package error

import (
    "bytespace.network/rerect/span"
)

// The error struct
// ----------------
type Error struct {
    Unit     CompUnit
    Position span.Span
    Message  string
}

func NewError(unit CompUnit, pos span.Span, msg string) Error {
    return Error {
        Unit: unit,
        Position: pos,
        Message: msg,
    }
}

// Error units (where did the error occour?)
// -----------------------------------------
type CompUnit string
const (
    FIO CompUnit = "FileIO"
    LEX CompUnit = "Lexer"
    PRS CompUnit = "Parser"
    BND CompUnit = "Binder"
    RNT CompUnit = "Runtime"
)
