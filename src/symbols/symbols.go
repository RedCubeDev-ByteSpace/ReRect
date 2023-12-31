package symbols

// Symbol base interface
// ---------------------
type Symbol interface {
    Type() SymbolType
    Name() string
}

// Types of sybols
// ---------------
type SymbolType string

const (
    ST_Package   SymbolType = "Package symbol"
    ST_Function  SymbolType = "Function symbol"
    ST_Global    SymbolType = "Global symbol"
    ST_Local     SymbolType = "Local symbol"
    ST_Parameter SymbolType = "Parameter symbol"
    ST_Type      SymbolType = "Type symbol"
    ST_Container SymbolType = "Container symbol"
    ST_Field     SymbolType = "Field symbol"
    ST_Instance  SymbolType = "Instance symbol"
    ST_Trait     SymbolType = "Trait symbol"
)
