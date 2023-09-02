// CompUnit - compunit.go
// ---------------------------------------------------------
// This packages holds all sorts of compilation related info
// ---------------------------------------------------------
package compunit

import "fmt"

// Source file struct
// ------------------
type SourceFile struct {
    Path    string
    Content string
}

// Register of source files
// ------------------------
var SourceFileRegister []SourceFile = make([]SourceFile, 0)

func RegisterSource(file string, content string) int {
    SourceFileRegister = append(SourceFileRegister, SourceFile{
        Path: file,
        Content: content,
    })

    fmt.Printf("Has registered file '%s': \n%s\n", file, content)

    return len(SourceFileRegister) - 1
}
