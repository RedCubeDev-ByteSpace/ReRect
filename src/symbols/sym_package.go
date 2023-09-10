package symbols

import "slices"

// Package variable symbol
// -----------------------
type PackageSymbol struct {
    Symbol

    PackName string
    Functions []*FunctionSymbol
    Globals []*GlobalSymbol
    Containers []*ContainerSymbol
    SymbolNames []string

    LoadedPackages map[string]*PackageSymbol
    IncludedPackages []string
}

func NewPackageSymbol(name string, funcs []*FunctionSymbol) *PackageSymbol {
    return &PackageSymbol{
        PackName: name,
        Functions: funcs,

        LoadedPackages: make(map[string]*PackageSymbol),
        IncludedPackages: make([]string, 0),
    }
}

func (sym *PackageSymbol) Name() string {
    return sym.PackName
}

func (sym *PackageSymbol) Type() SymbolType {
    return ST_Package
}


func (sym *PackageSymbol) TryRegisterContainer(cnt *ContainerSymbol) bool {
    // check if a symbol with this name already exists
    if slices.Contains(sym.SymbolNames, cnt.Name()) {
        return false
    }

    sym.Containers = append(sym.Containers, cnt)
    sym.SymbolNames = append(sym.SymbolNames, cnt.Name())
    return true
}

func (sym *PackageSymbol) TryRegisterFunction(fnc *FunctionSymbol) bool {
    // check if a symbol with this name already exists
    if slices.Contains(sym.SymbolNames, fnc.Name()) {
        return false
    }

    sym.Functions = append(sym.Functions, fnc)
    sym.SymbolNames = append(sym.SymbolNames, fnc.Name())
    return true
}

func (sym *PackageSymbol) TryRegisterGlobal(glb *GlobalSymbol) bool {
    // check if a symbol with this name already exists
    if slices.Contains(sym.SymbolNames, glb.Name()) {
        return false
    }

    sym.Globals = append(sym.Globals, glb)
    sym.SymbolNames = append(sym.SymbolNames, glb.Name())
    return true
}
