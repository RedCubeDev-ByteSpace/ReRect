package symbols

// Local variable symbol
// ---------------------
type LocalSymbol struct {
    VariableSymbol

    VarName string
    VarType *TypeSymbol
}

func NewLocalSymbol(name string, typ *TypeSymbol) *LocalSymbol {
    return &LocalSymbol{
        VarName: name,
        VarType: typ,
    }
}

func (sym *LocalSymbol) Name() string {
    return sym.VarName
}

func (sym *LocalSymbol) Type() SymbolType {
    return ST_Local
}
