package symbols

// Field variable symbol
// ---------------------
type FieldSymbol struct {
    VariableSymbol

    ParentContainer *ContainerSymbol

    FieldName string
    FieldType *TypeSymbol
}

func NewFieldSymbol(cnt *ContainerSymbol, name string, typ *TypeSymbol) *FieldSymbol {
    return &FieldSymbol{
        ParentContainer: cnt,
        FieldName: name,
        FieldType: typ,
    }
}

func (sym *FieldSymbol) Name() string {
    return sym.FieldName
}

func (sym *FieldSymbol) Type() SymbolType {
    return ST_Field
}

func (sym *FieldSymbol) VarType() *TypeSymbol {
    return sym.FieldType
}
