package symbols

// Container variable symbol
// ---------------------
type ContainerSymbol struct {
    Symbol

    ParentPackage *PackageSymbol

    ContainerName string
    ContainerType *TypeSymbol

    Constructor *FunctionSymbol

    Symbols []string
    Fields []*FieldSymbol
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
