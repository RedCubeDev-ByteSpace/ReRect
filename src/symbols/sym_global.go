package symbols

// Global variable symbol
// ---------------------
type GlobalSymbol struct {
    VariableSymbol

    VarName string
    VarType *TypeSymbol
}

func NewGlobalSymbol(name string, typ *TypeSymbol) *GlobalSymbol {
    return &GlobalSymbol{
        VarName: name,
        VarType: typ,
    }
}

func (sym *GlobalSymbol) Name() string {
    return sym.VarName
}

func (sym *GlobalSymbol) Type() SymbolType {
    return ST_Global
}
