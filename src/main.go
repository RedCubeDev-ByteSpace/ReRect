package main

import (
	"fmt"
	"slices"

	"bytespace.network/rerect/binder"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	packageprocessor "bytespace.network/rerect/package_processor"
	"bytespace.network/rerect/parser"
	"bytespace.network/rerect/syntaxnodes"
)

func main() {
    
    // Lexing
    // ------
    tokens := lexer.LexFile("./tests/print.rr")
    
    fmt.Println("\nLexer output:")
    for _, v := range tokens {
        fmt.Printf("%s, %s, '%s'\n", v.Type, v.Position.Format(), v.Buffer)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Parsing
    // -------
    members := parser.Parse(tokens)

    fmt.Println("\nParser output:")
    for _, v := range members {
        fmt.Printf("%s\n", v.Type())
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Package processing
    // ------------------
    packageprocessor.Init()
    packs, mems := packageprocessor.Process([][]syntaxnodes.MemberNode{members})

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

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Binding
    // -------

    // First: Index all functions
    for i, _ := range mems {
        binder.IndexFunctions(packs[i], mems[i])
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

}
