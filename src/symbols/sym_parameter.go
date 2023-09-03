package symbols

// Parameter variable symbol
// ---------------------
type ParameterSymbol struct {
    VariableSymbol

    VarName string
    VarIdx  int
    VarType *TypeSymbol
}

func NewParameterSymbol(name string, idx int, typ *TypeSymbol) *ParameterSymbol {
    return &ParameterSymbol{
        VarName: name,
        VarIdx: idx,
        VarType: typ,
    }
}

func (sym *ParameterSymbol) Name() string {
    return sym.VarName
}

func (sym *ParameterSymbol) Type() SymbolType {
    return ST_Parameter
}
