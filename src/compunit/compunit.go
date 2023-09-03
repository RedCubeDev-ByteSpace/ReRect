// CompUnit - compunit.go
// ---------------------------------------------------------
// This packages holds all sorts of compilation related info
// ---------------------------------------------------------
package compunit

import (
	"fmt"

	"bytespace.network/rerect/symbols"
)

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

// Register of global data types
// -----------------------------
var GlobalDataTypeRegister map[string]*symbols.TypeSymbol = map[string]*symbols.TypeSymbol {
    "void": symbols.NewTypeSymbol("void", make([]*symbols.TypeSymbol, 0), symbols.NONE), // nothin
    
    "long": symbols.NewTypeSymbol("long", make([]*symbols.TypeSymbol, 0), symbols.INT), // 64 bit int
    "int" : symbols.NewTypeSymbol("int" , make([]*symbols.TypeSymbol, 0), symbols.INT), // 32 bit int
    "word": symbols.NewTypeSymbol("word", make([]*symbols.TypeSymbol, 0), symbols.INT), // 16 bit int
    "byte": symbols.NewTypeSymbol("byte", make([]*symbols.TypeSymbol, 0), symbols.INT), // 8  bit int

    "float" : symbols.NewTypeSymbol("float" , make([]*symbols.TypeSymbol, 0), symbols.FLOAT), // 32 bit float
    "double": symbols.NewTypeSymbol("dobule", make([]*symbols.TypeSymbol, 0), symbols.FLOAT), // 64 bit float
    
    "string": symbols.NewTypeSymbol("string", make([]*symbols.TypeSymbol, 0), symbols.NONE), // string
}

// Register of known packages
// --------------------------
var PackagesRegister map[string]*symbols.PackageSymbol = make(map[string]*symbols.PackageSymbol)

// A few helpers for this one
func GetPackage(name string) *symbols.PackageSymbol {
    pck, ok := PackagesRegister[name]

    if !ok {
        return nil
    }

    return pck
}

func CreatePackage(name string) *symbols.PackageSymbol {
    _, ok := PackagesRegister[name]

    if ok {
        return nil
    }

    PackagesRegister[name] = symbols.NewPackageSymbol(name, make([]*symbols.FunctionSymbol, 0))
    return PackagesRegister[name]
}

func GetPackageAtAllCosts(name string) *symbols.PackageSymbol {
    pck := GetPackage(name)

    if pck != nil {
        return pck
    }

    return CreatePackage(name)
}

