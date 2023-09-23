package symbols

// Function variable symbol
// ---------------------
type FunctionSymbol struct {
	Symbol

	FunctionKind FunctionType
	MethodKind   MethodType
	MethodSource *TypeSymbol

    SourceTrait *TraitSymbol

    // if this is a method implemented by a trait -> mark it as needing to be redirected and add a reference
    NeedsVirtualCallToTrait bool
    TraitSourceMethod *FunctionSymbol 

    // if this is a trait method implemented by the container -> mark it as needing to be redirected to the container
    NeedsVirtualCallToContainer bool

	ParentPackage *PackageSymbol

	FuncName   string
	ReturnType *TypeSymbol

	IsVMFunction    bool
	FunctionPointer VMFPtr
	MethodPointer   VMMPtr

	Parameters []*ParameterSymbol
}

type VMFPtr func([]interface{}) interface{}
type VMMPtr func(interface{}, []interface{}) interface{}

func NewFunctionSymbol(pck *PackageSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol) *FunctionSymbol {
	return &FunctionSymbol{
		FunctionKind:  FT_FUNC,
		ParentPackage: pck,

		FuncName:   name,
		ReturnType: typ,
		Parameters: params,

		IsVMFunction: false,
	}
}

func NewMethodSymbol(pck *PackageSymbol, src *TypeSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol) *FunctionSymbol {
	return &FunctionSymbol{
		FunctionKind: FT_METH,
		MethodKind:   MT_STRICT,
		MethodSource: src,

		ParentPackage: pck,

		FuncName:   name,
		ReturnType: typ,
		Parameters: params,

		IsVMFunction: false,
	}
}

func NewVMFunctionSymbol(pck *PackageSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol, ptr VMFPtr) *FunctionSymbol {
	return &FunctionSymbol{
		FunctionKind:  FT_FUNC,
		ParentPackage: pck,

		FuncName:   name,
		ReturnType: typ,
		Parameters: params,

		IsVMFunction:    true,
		FunctionPointer: ptr,
	}
}

func NewVMMethodSymbol(pck *PackageSymbol, meth MethodType, src *TypeSymbol, name string, typ *TypeSymbol, params []*ParameterSymbol, ptr VMMPtr) *FunctionSymbol {
	return &FunctionSymbol{
		FunctionKind: FT_METH,
		MethodKind:   meth,
		MethodSource: src,

		ParentPackage: pck,

		FuncName:   name,
		ReturnType: typ,
		Parameters: params,

		IsVMFunction:  true,
		MethodPointer: ptr,
	}
}

func (sym *FunctionSymbol) Name() string {
	return sym.FuncName
}

func (sym *FunctionSymbol) Type() SymbolType {
	return ST_Function
}

type FunctionType string

const (
	FT_FUNC FunctionType = "Function"
	FT_METH FunctionType = "Method"
)

type MethodType string

const (
	MT_STRICT MethodType = "Strict" // Just a single, specific type
	MT_GROUP  MethodType = "Group"  // All types of a single group (all ints, all floats, all arrays, or all containers)
	MT_ALL    MethodType = "All"    // Literally all types
)
