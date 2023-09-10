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
    Globals []*symbols.GlobalSymbol
}

func compFailed() *CompilationResult {
    return &CompilationResult{
        Ok: false,
    }
}

const DBG = false

func Compile(srcFiles []string) *CompilationResult {
    // Lexing
    // ------
    fileTokens := [][]lexer.Token{}

    for _, file := range srcFiles {
        tokens := lexer.LexFile(file)

        if DBG {
            fmt.Println("\nLexer output:")
            for _, v := range tokens {
                fmt.Printf("%s, %s, '%s'\n", v.Type, v.Position.Format(), v.Buffer)
            }
        }

        fileTokens = append(fileTokens, tokens)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Parsing
    // -------
    fileMembers := [][]syntaxnodes.MemberNode{}

    for _, tokens := range fileTokens {
        members := parser.Parse(tokens)

        if DBG {
            fmt.Println("\nParser output:")
            for _, v := range members {
                fmt.Printf("%s\n", v.Type())
            }
        }

        fileMembers = append(fileMembers, members)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Package processing
    // ------------------
    packageprocessor.Init()
    files := packageprocessor.Process(fileMembers)

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

    // Even Firsterer: Index all container datatypes (this NEEDS to be done first!!!! otherwise stuff cant be linked correctly!)
    for _, file := range files {
        binder.IndexContainerTypes(file)
    }

    // A little Firster: Index all container fields and methods
    for _, file := range files {
        binder.IndexContainerContents(file)
    }

    // First: Index all functions and globals
    for _, file := range files {
        binder.IndexFunctions(file)
        binder.IndexGlobals(file)
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

    if error.HasErrors() {
        error.Output()
        return compFailed()
    }

    // Bring the compilation result into a usable format
    // -------------------------------------------------
    res := &CompilationResult{
        Ok: true,
        Functions: make(map[*symbols.FunctionSymbol]*boundnodes.BoundBlockStatementNode),
        Globals: make([]*symbols.GlobalSymbol, 0),
    }

    for _, file := range files {
        for k, v := range file.FunctionBodies {
            res.Functions[k] = v.(*boundnodes.BoundBlockStatementNode)
        }

        for _, v := range file.Globals {
            res.Globals = append(res.Globals, v)
        }
    }

    return res
}
