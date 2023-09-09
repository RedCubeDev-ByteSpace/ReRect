package evalobjects

import "bytespace.network/rerect/symbols"

// Implementations for the array type
// ----------------------------------
type ArrayInstance struct {
    Type *symbols.TypeSymbol
    Elements []interface{}
}
