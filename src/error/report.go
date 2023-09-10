// Error - report.go
// ---------------------------------------------------------------------
// Contains the utilities for actually reporting errors
// ---------------------------------------------------------------------
package error

import (
	"fmt"
	"os"
	"strings"

	"bytespace.network/rerect/compunit"
)

// ANSI color constants
// --------------------
const (
    RED = "\033[31m"
    GRN = "\033[32m"
    YLW = "\033[33m"
    CYN = "\033[36m"
    RST = "\033[0m"
)

// Error collection
// ----------------
var errors []Error

// Add an error to the collection
// ------------------------------
func Report(err Error) {
    errors = append(errors, err)

    // if error occoured in runtime -> this is bad :)
    if err.Unit == RNT {
        Output()
        os.Exit(-1)
    }
} 

// Output all reported errors
// --------------------------
func Output() {
    for _, err := range errors {
        fmt.Print(RED)
   
        if !err.Position.Internal {
            // figure out where the error happened
            line, col := err.Position.GetLineAndCol()

            // get the text from that line
            errline := strings.Split(compunit.SourceFileRegister[err.Position.File].Content, "\n")[line - 1]
            errline = strings.Replace(errline, "\t", " ", -1)

            // calculate length of error underline
            underlineLen := err.Position.ToIdx - err.Position.FromIdx
            if underlineLen == 0 {
                underlineLen = 1
            }

            fmt.Printf("[%s][L:%d, C:%d]: %s\n", err.Unit, line, col, err.Message)
            fmt.Print(RST)
            fmt.Println(errline)
            fmt.Printf("%s%s%s%s\n", RED, strings.Repeat(" ", col-1), strings.Repeat("^", underlineLen), RST)
        } else {
            fmt.Printf("[%s][Internal]: %s\n", err.Unit, err.Message)
            fmt.Print(RST)
        }
        fmt.Println()
    }
}

// Are there errors?
// -----------------
func HasErrors() bool {
    return len(errors) > 0
}
