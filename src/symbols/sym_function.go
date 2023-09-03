package symbols

// Function variable symbol
// ---------------------
type FunctionSymbol struct {
    VariableSymbol

    FuncName string
    ReturnType *TypeSymbol

    IsVMFunction bool
    Pointer interface{}

    Parameters []*ParameterSymbol
}

func NewFunctionSymbol(name string, typ *TypeSymbol, params []*ParameterSymbol) *FunctionSymbol {
    return &FunctionSymbol{
        FuncName: name,
        ReturnType: typ,
        Parameters: params,
        IsVMFunction: false,
    }
}

func NewVMFunctionSymbol(name string, typ *TypeSymbol, params []*ParameterSymbol, ptr interface{}) *FunctionSymbol {
    return &FunctionSymbol{
        FuncName: name,
        ReturnType: typ,
        Parameters: params,

        IsVMFunction: true,
        Pointer: ptr,
    }
}

func (sym *FunctionSymbol) Name() string {
    return sym.FuncName
}

func (sym *FunctionSymbol) Type() SymbolType {
    return ST_Function
}
