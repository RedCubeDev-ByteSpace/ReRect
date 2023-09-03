package symbols

// Symbol base interface
// ---------------------
type Symbol interface {
    Type() SymbolType
    Name() string
    Equal(*Symbol) bool
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
)
