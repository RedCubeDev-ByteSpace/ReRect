package symbols

// Container symbol
// ----------------
type ContainerSymbol struct {
    Symbol

    ParentPackage *PackageSymbol
    Traits []*TraitSymbol

    ContainerName string
    ContainerType *TypeSymbol

    Constructor *FunctionSymbol

    Symbols []string
    Fields []*FieldSymbol
    Methods []*FunctionSymbol
}

func NewContainerSymbol(pck *PackageSymbol, name string, typ *TypeSymbol) *ContainerSymbol {
    cnt := &ContainerSymbol{
        ParentPackage: pck,
        ContainerName: name,
        ContainerType: typ,

        // Fields will be filled in later
        Fields: make([]*FieldSymbol, 0),

        // same here
        Symbols: make([]string, 0),

        // same here
        Traits: make([]*TraitSymbol, 0),

        // same here
        Methods: make([]*FunctionSymbol, 0),
    }

    // link the given type symbol to this container
    typ.Container = cnt

    // ok we don
    return cnt
}

func (sym *ContainerSymbol) Name() string {
    return sym.ContainerName
}

func (sym *ContainerSymbol) Type() SymbolType {
    return ST_Container
}

func (sym *ContainerSymbol) VarType() *TypeSymbol {
    return sym.ContainerType
}
