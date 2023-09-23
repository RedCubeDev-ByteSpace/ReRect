package symbols

// Trait symbol
// ------------
type TraitSymbol struct {
    Symbol

    ParentPackage *PackageSymbol

    TraitName string
    TraitType *TypeSymbol

    Symbols []string
    Fields []*FieldSymbol
    Methods []*FunctionSymbol
}

func NewTraitSymbol(pck *PackageSymbol, name string, typ *TypeSymbol) *TraitSymbol {
    cnt := &TraitSymbol{
        ParentPackage: pck,
        TraitName: name,
        TraitType: typ,

        // Fields will be filled in later
        Fields: make([]*FieldSymbol, 0),

        // same here
        Symbols: make([]string, 0),

        // even samer here
        Methods: make([]*FunctionSymbol, 0),
    }

    // link the given type symbol to this container
    typ.Trait = cnt

    // ok we don
    return cnt
}

func (sym *TraitSymbol) Name() string {
    return sym.TraitName
}

func (sym *TraitSymbol) Type() SymbolType {
    return ST_Trait
}

func (sym *TraitSymbol) VarType() *TypeSymbol {
    return sym.TraitType
}
