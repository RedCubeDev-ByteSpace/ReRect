package symbols

// Global variable symbol
// ---------------------
type GlobalSymbol struct {
    VariableSymbol

    GlobalName string
    GlobalType *TypeSymbol
}

func NewGlobalSymbol(name string, typ *TypeSymbol) *GlobalSymbol {
    return &GlobalSymbol{
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
