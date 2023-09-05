package compctl

import (
	"fmt"
	"slices"

	"bytespace.network/rerect/binder"
	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/lowerer"
	"bytespace.network/rerect/package_processor"
	"bytespace.network/rerect/parser"
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

type CompilationResult struct {
    Ok bool
    Functions map[*symbols.FunctionSymbol]*boundnodes.BoundBlockStatementNode
}

func compFailed() *CompilationResult {
    return &CompilationResult{
        Ok: false,
    }
}

const DBG = false

func Compile(file string) *CompilationResult {
    // Lexing
    // ------
    tokens := lexer.LexFile(file)
  
    if DBG {
        fmt.Println("\nLexer output:")
        for _, v := range tokens {
            fmt.Printf("%s, %s, '%s'\n", v.Type, v.Position.Format(), v.Buffer)
        }
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Parsing
    // -------
    members := parser.Parse(tokens)

    if DBG {
        fmt.Println("\nParser output:")
        for _, v := range members {
            fmt.Printf("%s\n", v.Type())
        }
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Package processing
    // ------------------
    packageprocessor.Init()
    files := packageprocessor.Process([][]syntaxnodes.MemberNode{members})

    if DBG {
        fmt.Println("\nPackage processor output:")
        for _, v := range compunit.PackagesRegister {
            fmt.Printf("Package: %s\n", v.PackName)
            fmt.Println(" Loads:")

            for k, _ := range v.LoadedPackages {
                fmt.Printf("  - %s ", k)

                if slices.Contains(v.IncludedPackages, k) {
                    fmt.Println("(included)")
                } else {
                    fmt.Println()
                }
            }
        }
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Binding
    // -------

    // First: Index all functions
    for _, file := range files {
        binder.IndexFunctions(file)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Second: bind all function bodies
    for _, file := range files {
        binder.BindFunctions(file)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Lowering
    // --------
    for _, file := range files {
        lowerer.Lower(file)
    }

    // Bring the compilation result into a usable format
    // -------------------------------------------------
    res := &CompilationResult{
        Ok: true,
        Functions: make(map[*symbols.FunctionSymbol]*boundnodes.BoundBlockStatementNode),
    }

    for _, file := range files {
        for k, v := range file.FunctionBodies {
            res.Functions[k] = v.(*boundnodes.BoundBlockStatementNode)
        }
    }

    return res
}
