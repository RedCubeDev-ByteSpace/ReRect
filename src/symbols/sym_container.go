package symbols

// Container variable symbol
// ---------------------
type ContainerSymbol struct {
    Symbol

    ParentPackage *PackageSymbol

    ContainerName string
    ContainerType *TypeSymbol

    Fields []*FieldSymbol
}

func NewContainerSymbol(pck *PackageSymbol, name string, typ *TypeSymbol) *ContainerSymbol {
    return &ContainerSymbol{
        ParentPackage: pck,
        ContainerName: name,
        ContainerType: typ,

        // Fields will be filled in later
        Fields: make([]*FieldSymbol, 0),
    }
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
