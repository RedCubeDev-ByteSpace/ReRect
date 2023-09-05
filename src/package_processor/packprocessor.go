// Package processor - packprocessor.go
// --------------------------------------------------------
// What packages exist and how are they linked up?
// --------------------------------------------------------
package packageprocessor

import (
	"bytespace.network/rerect/boundnodes"
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	gopackages "bytespace.network/rerect/go_packages"
	"bytespace.network/rerect/symbols"
	"bytespace.network/rerect/syntaxnodes"
)

type CompilationFile struct {
    Package *symbols.PackageSymbol
    Members []syntaxnodes.MemberNode

    Functions []*symbols.FunctionSymbol
    FunctionBodiesSrc map[*symbols.FunctionSymbol]syntaxnodes.StatementNode
    FunctionBodies    map[*symbols.FunctionSymbol]boundnodes.BoundStatementNode
}

func Init() {
    // load all native packages
    gopackages.Load()
}

func Process(mems [][]syntaxnodes.MemberNode) []*CompilationFile {
    files := []*CompilationFile{}

    // Register all names first
    for _, mem := range mems {
        p := register(mem)
        //fmt.Println(p.PackName)

        files = append(files, &CompilationFile{
            Package: p,
            Members: mem,

            // create containers to be filled in by the binder later
            Functions: []*symbols.FunctionSymbol{},
            FunctionBodiesSrc: make(map[*symbols.FunctionSymbol]syntaxnodes.StatementNode),
            FunctionBodies: make(map[*symbols.FunctionSymbol]boundnodes.BoundStatementNode),
        })
    }
    
    // Link up the packages
    for _, file := range files {
        link(file.Package, file.Members)
    }

    return files
}

func register(mem []syntaxnodes.MemberNode) *symbols.PackageSymbol {
    packageName := "main"

    // search through all members
    for _, nd := range mem {
        // we're only looking for package names
        if nd.Type() != syntaxnodes.NT_Package {
            continue
        } 

        // get the package name that was set
        node := nd.(*syntaxnodes.PackageNode)
        packageName = node.PackageName.Buffer
    }

    // get or create the package of that name
    pack := compunit.GetPackageAtAllCosts(packageName)
    return pack
}

func link(pck *symbols.PackageSymbol, mem []syntaxnodes.MemberNode) {

    // search through all members
    for _, nd := range mem {
        // we're only looking for load statements
        if nd.Type() != syntaxnodes.NT_Load {
            continue
        } 

        // get the package name that we need to load
        node := nd.(*syntaxnodes.LoadNode)
        packageName := node.Library.Buffer

        // look the package up and add a reference
        ref := compunit.GetPackage(packageName)

        // lookup failed
        if ref == nil {
            error.Report(error.NewError(error.PCK, node.Position(), "Could not find package '%s'!", packageName))
            continue
        }

        // if the lookup succeeded -> add the ref
        pck.LoadedPackages[packageName] = ref

        // is the package included? -> if so: add it to the list
        if node.Included {
            pck.IncludedPackages = append(pck.IncludedPackages, packageName)
        }
    }
}  
