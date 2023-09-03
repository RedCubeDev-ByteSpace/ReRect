package symbols

// Local variable symbol
// ---------------------
type LocalSymbol struct {
    VariableSymbol

    LocalName string
    LocalType *TypeSymbol
}

func NewLocalSymbol(name string, typ *TypeSymbol) *LocalSymbol {
    return &LocalSymbol{
        LocalName: name,
        LocalType: typ,
    }
}

func (sym *LocalSymbol) Name() string {
    return sym.LocalName
}

func (sym *LocalSymbol) Type() SymbolType {
    return ST_Local
}

func (sym *LocalSymbol) VarType() *TypeSymbol {
    return sym.LocalType
}
