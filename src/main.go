package main

import (
	"fmt"
	"slices"

	"bytespace.network/rerect/binder"
	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/lowerer"
	packageprocessor "bytespace.network/rerect/package_processor"
	"bytespace.network/rerect/parser"
	"bytespace.network/rerect/symbols"
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
    type PackageUnit struct {
        Package *symbols.PackageSymbol
        FunctionSymbols []*symbols.FunctionSymbol
        SourceFunctionBodies []syntaxnodes.StatementNode
        FunctionBodies []boundnodes.BoundStatementNode
    }
    
    pus := []PackageUnit{}

    // First: Index all functions
    for i, _ := range mems {
        syms, srcbodies := binder.IndexFunctions(packs[i], mems[i])
        pus = append(pus, PackageUnit{
            Package: packs[i],
            FunctionSymbols: syms,
            SourceFunctionBodies: srcbodies,
        })
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Second: bind all function bodies
    for i, v := range pus {
        bodies := binder.BindFunctions(v.Package, v.FunctionSymbols, v.SourceFunctionBodies)
        pus[i].FunctionBodies = bodies
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Lowering
    // --------
    fmt.Println("\nLowerer output:")

    for _, p := range pus {
        fmt.Printf("Package unit %s\n", p.Package.Name())

        for i, _ := range p.FunctionSymbols {
            fmt.Printf("function %s:\n",p.FunctionSymbols[i].FuncName)
            p.FunctionBodies[i] = lowerer.Lower(p.FunctionBodies[i])
            
            for _, v := range p.FunctionBodies[i].(*boundnodes.BoundBlockStatementNode).Statements {
                fmt.Printf("  %s\n", v.Type())
            }
        }
    }

}
