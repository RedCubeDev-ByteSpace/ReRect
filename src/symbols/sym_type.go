package symbols

// Type symbol struct
// ------------------
type TypeSymbol struct {
    Symbol

    TypeName string
    SubTypes []*TypeSymbol   

    TypeGroup TypeGroupType
    TypeSize int

    Default interface{}
}

func NewTypeSymbol(name string, subtypes []*TypeSymbol, grp TypeGroupType, sz int, def interface{}) *TypeSymbol {
    return &TypeSymbol {
        TypeName: name,
        SubTypes: subtypes,
        TypeGroup: grp,
        TypeSize: sz,
        Default: def,
    }
}

func  (typ *TypeSymbol) Name() string {
    return typ.TypeName
}

func (typ *TypeSymbol) Type() SymbolType {
    return ST_Type
}

func (t1 *TypeSymbol) Equal(t2 *TypeSymbol) bool {
    // Names dont match? -> they're DEFINITELY not the same lol
    if t1.Name() != t2.Name() {
        return false
    }

    // Do both types have the same size?
    if t1.TypeSize != t2.TypeSize {
        return false
    }

    // Do both types have the same amount of subtypes?
    if len(t1.SubTypes) != len(t2.SubTypes) {
        return false
    }

    // Are the subtypes equal?
    for i := range t1.SubTypes {
        // make sure the types are equal
        if !t1.SubTypes[i].Equal(t2.SubTypes[i]) {
            return false
        }
    }

    // damn okay ig theyre identical
    return true
} 

// Types of data types
// -------------------
type TypeGroupType string;
const (
    NONE  TypeGroupType = "No group"
    INT   TypeGroupType = "Integer type"
    FLOAT TypeGroupType = "Floating point type"
    ARR   TypeGroupType = "Array type"
    CONT  TypeGroupType = "Container type"
)
