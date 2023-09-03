package symbols

// Parameter variable symbol
// ---------------------
type ParameterSymbol struct {
    VariableSymbol

    ParameterName string
    ParameterIdx  int
    ParameterType *TypeSymbol
}

func NewParameterSymbol(name string, idx int, typ *TypeSymbol) *ParameterSymbol {
    return &ParameterSymbol{
        ParameterName: name,
        ParameterIdx: idx,
        ParameterType: typ,
    }
}

func (sym *ParameterSymbol) Name() string {
    return sym.ParameterName
}

func (sym *ParameterSymbol) Type() SymbolType {
    return ST_Parameter
}

func (sym *ParameterSymbol) VarType() *TypeSymbol {
    return sym.ParameterType
}
