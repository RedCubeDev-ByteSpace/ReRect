package symbols

// Type symbol struct
// ------------------
type TypeSymbol struct {
    Symbol

    TypeName string
    SubTypes []*TypeSymbol   

    TypeGroup TypeGroupType
    TypeSize int
}

func NewTypeSymbol(name string, subtypes []*TypeSymbol, grp TypeGroupType, sz int) *TypeSymbol {
    return &TypeSymbol {
        TypeName: name,
        SubTypes: subtypes,
        TypeGroup: grp,
        TypeSize: sz,
    }
}

func  (typ *TypeSymbol) Name() string {
    return typ.TypeName
}

func (typ *TypeSymbol) Type() SymbolType {
    return ST_Type
}

func (t1 *TypeSymbol) Equal(t2 *TypeSymbol) bool {
    return t1.Name() == t2.Name()
} 

// Types of data types
// -------------------
type TypeGroupType string;
const (
    NONE  TypeGroupType = "No group"
    INT   TypeGroupType = "Integer type"
    FLOAT TypeGroupType = "Floating point type"
)
