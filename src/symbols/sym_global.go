package symbols

// Global variable symbol
// ---------------------
type GlobalSymbol struct {
    VariableSymbol

    ParentPackage *PackageSymbol

    GlobalName string
    GlobalType *TypeSymbol
}

func NewGlobalSymbol(pck *PackageSymbol, name string, typ *TypeSymbol) *GlobalSymbol {
    return &GlobalSymbol{
        ParentPackage: pck,
        GlobalName: name,
        GlobalType: typ,
    }
}

func (sym *GlobalSymbol) Name() string {
    return sym.GlobalName
}

func (sym *GlobalSymbol) Type() SymbolType {
    return ST_Global
}

func (sym *GlobalSymbol) VarType() *TypeSymbol {
    return sym.GlobalType
}
