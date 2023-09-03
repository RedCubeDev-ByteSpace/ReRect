package symbols

// Type symbol struct
// ------------------
type TypeSymbol struct {
    Symbol

    TypeName string
    SubTypes []*TypeSymbol   

    TypeGroup TypeGroupType
}

func NewTypeSymbol(name string, subtypes []*TypeSymbol, grp TypeGroupType) *TypeSymbol {
    return &TypeSymbol {
        TypeName: name,
        SubTypes: subtypes,
        TypeGroup: grp,
    }
}

func  (typ *TypeSymbol) Name() string {
    return typ.TypeName
}

func (typ *TypeSymbol) Type() SymbolType {
    return ST_Type
}

// Types of data types
// -------------------
type TypeGroupType string;
const (
    NONE  TypeGroupType = "No group"
    INT   TypeGroupType = "Integer type"
    FLOAT TypeGroupType = "Floating point type"
)
