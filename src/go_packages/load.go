// Go Packages - loader.go
// --------------------------------------------------------
// Loads all them packages written in go
// --------------------------------------------------------
package gopackages

import (
	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/symbols"
)

func Load() {
    // Load the sys package
    LoadSys()
}

// A few helpers
// -------------
func registerPackage(name string) {
    pck := compunit.GetPackage(name)
    if pck != nil {
        error.Report(error.NewError(error.GOP, span.Internal(), "Unable to register package '%s'! A package with that name already exists!", name))
    }

    compunit.CreatePackage(name)
}

func registerFunction(pack string, fnc *symbols.FunctionSymbol) {
    pck := compunit.GetPackage(pack)
    
    if pck == nil {
        error.Report(error.NewError(error.GOP, span.Internal(), "Unable to register function '%s' in package '%s'! No package called '%s' could be found!", fnc.FuncName, pack, pack))
    }

    pck.Functions = append(pck.Functions, fnc)
}
