package evalobjects

import "bytespace.network/rerect/symbols"

// Container instance struct
// -------------------------
type ContainerInstance struct {
    Type *symbols.TypeSymbol
    Fields map[string]interface{}
}
