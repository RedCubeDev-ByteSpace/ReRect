package symbols

// Instance variable symbol
// ------------------------
type InstanceSymbol struct {
    VariableSymbol

    InstanceType *TypeSymbol
}

func NewInstanceSymbol(typ *TypeSymbol) *InstanceSymbol {
    return &InstanceSymbol{
        InstanceType: typ,
    }
}

func (sym *InstanceSymbol) Name() string {
    return "this"
}

func (sym *InstanceSymbol) Type() SymbolType {
    return ST_Instance
}

func (sym *InstanceSymbol) VarType() *TypeSymbol {
    return sym.InstanceType
}
