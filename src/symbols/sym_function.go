package symbols

// Function variable symbol
// ---------------------
type FunctionSymbol struct {
    VariableSymbol

    ParentPackage *PackageSymbol

    FuncName string
    ReturnType *TypeSymbol

    IsVMFunction bool
    Pointer VMFPtr

    Parameters []*ParameterSymbol
}

type VMFPtr func([]interface{}) interface{}

func NewFunctionSymbol(pck *PackageSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol) *FunctionSymbol {
    return &FunctionSymbol{
        ParentPackage: pck,

        FuncName: name,
        ReturnType: typ,
        Parameters: params,

        IsVMFunction: false,
    }
}

func NewVMFunctionSymbol(pck *PackageSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol, ptr VMFPtr) *FunctionSymbol {
    return &FunctionSymbol{
        ParentPackage: pck,

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
