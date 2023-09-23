package symbols

// Field variable symbol
// ---------------------
type FieldSymbol struct {
    VariableSymbol

    ParentContainer *ContainerSymbol
    HasParentContainer bool

    ParentTrait *TraitSymbol
    HasParentTrait bool

    FieldName string
    FieldType *TypeSymbol
}

func NewFieldSymbol(cnt *ContainerSymbol, name string, typ *TypeSymbol) *FieldSymbol {
    return &FieldSymbol{
        ParentContainer: cnt,
        HasParentContainer: true,
        FieldName: name,
        FieldType: typ,
    }
}

func NewTraitFieldSymbol(trt *TraitSymbol, name string, typ *TypeSymbol) *FieldSymbol {
    return &FieldSymbol{
        ParentTrait: trt,
        HasParentTrait: true,
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
