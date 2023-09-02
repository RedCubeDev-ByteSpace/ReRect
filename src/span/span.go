// Span - span.go
// ---------------------------------------------------------------------
// Contains utils for span and text things
// ---------------------------------------------------------------------
package span

import (
	"fmt"

	"bytespace.network/rerect/compunit"
)

type Span struct {
    File    int // File ID
    FromIdx int // Char start
    ToIdx   int // Char end
}

// Format span as a string
// -----------------------
func (s Span) Format() string {
    return fmt.Sprintf("%d, %d", s.FromIdx, s.ToIdx)
}

// Find span inbetween two spans
// -----------------------------
func (s1 Span) SpanBetween(s2 Span) Span {
    if s1.File != s2.File {
        return s1
    }

    spanBetween := s1

    // MIN(FromIdx)
    if s2.FromIdx < spanBetween.FromIdx {
        spanBetween.FromIdx = s2.FromIdx
    }

    // MAX(ToIdx)
    if s2.ToIdx > spanBetween.ToIdx {
        spanBetween.ToIdx = s2.ToIdx
    }

    return spanBetween
}

// Convert index values into line and column numbers
// -------------------------------------------------
func (s Span) GetLineAndCol() (line int, col int) {
    // get the contents of the file the span is in
    content := []rune(compunit.SourceFileRegister[s.File].Content)

    line = 1
    col = 1

    // count newlines up to the span
    for i := 0; i < s.FromIdx; i++ {
        if content[i] == '\n' {
            col = 0
            line++
        }

        col++
    }

    return
}
