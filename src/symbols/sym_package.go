package symbols


// Package variable symbol
// -----------------------
type PackageSymbol struct {
    VariableSymbol

    PackName string
    Functions []*FunctionSymbol
    Globals []*GlobalSymbol

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

func (sym *PackageSymbol) TryRegisterFunction(fnc *FunctionSymbol) bool {
    // check if a function with this name already exists
    for _, v := range sym.Functions {
        if v.FuncName == fnc.FuncName {
            // function name is already taken!

            return false
        }
    }

    sym.Functions = append(sym.Functions, fnc)
    return true
}

func (sym *PackageSymbol) TryRegisterGlobal(glb *GlobalSymbol) bool {
    // check if a global with this name already exists
    for _, v := range sym.Globals {
        if v.GlobalName == glb.GlobalName {
            // Global name is already taken!

            return false
        }
    }

    sym.Globals = append(sym.Globals, glb)
    return true
}
