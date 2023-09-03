package symbols

// Variable symbol base
// --------------------
type VariableSymbol interface {
    Symbol
    VarType() *TypeSymbol
} 
