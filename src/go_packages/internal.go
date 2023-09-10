package gopackages

import (
	"os"
	"reflect"

	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/error"
	evalobjects "bytespace.network/rerect/eval_objects"
	"bytespace.network/rerect/span"
	"bytespace.network/rerect/symbols"
)

// Some very cool internal functions
func LoadInternal() {
    pack := registerPackage("internal")
    
    // create a dummy array type symbol
    arr := symbols.NewTypeSymbol("array", []*symbols.TypeSymbol{}, symbols.ARR, 0, nil)

    // Array methods
    registerFunction("internal", symbols.NewVMMethodSymbol(pack, symbols.MT_GROUP , arr, "Length", compunit.GlobalDataTypeRegister["int"] , []*symbols.ParameterSymbol{}, Array_Length))
    registerFunction("internal", symbols.NewVMMethodSymbol(pack, symbols.MT_GROUP , arr, "Push"  , compunit.GlobalDataTypeRegister["void"], []*symbols.ParameterSymbol{symbols.NewParameterSymbol("Element", 0, compunit.GlobalDataTypeRegister["any"]) }, Array_Push))
    registerFunction("internal", symbols.NewVMMethodSymbol(pack, symbols.MT_GROUP , arr, "Pop"   , compunit.GlobalDataTypeRegister["any"] , []*symbols.ParameterSymbol{}, Array_Pop))

    // String methods
    registerFunction("internal", symbols.NewVMMethodSymbol(pack, symbols.MT_STRICT, compunit.GlobalDataTypeRegister["string"], "Length", compunit.GlobalDataTypeRegister["int"], []*symbols.ParameterSymbol{}, String_Length))

    // Global functions
    registerFunction("internal", symbols.NewVMFunctionSymbol(pack, "die", compunit.GlobalDataTypeRegister["void"], []*symbols.ParameterSymbol{symbols.NewParameterSymbol("exitcode", 0, compunit.GlobalDataTypeRegister["int"])}, Die))
}

func String_Length(instance any, args []any) any {
    // make sure the instance isnt null
    if instance == nil {
        return 0
    }

    // otherwise -> return the string length
    return int32(len(instance.(string)))
}

func Array_Length(instance any, args []any) any {
    // make sure the instance isnt null
    if instance == nil {
        return 0
    }

    // otherwise -> return the array length
    return int32(len(instance.(*evalobjects.ArrayInstance).Elements))
}

func Array_Push(instance any, args []any) any {
    // make sure the instance isnt null
    if instance == nil {
        return 0
    }

    // read out args
    arr := instance.(*evalobjects.ArrayInstance)
    elem := args[0]

    // make sure this arg is the correct type
    elem, ok := evalobjects.EvalConversion(elem, arr.Type.SubTypes[0])
    
    if !ok {
        typ := "any"
        if elem != nil {
            typ = reflect.TypeOf(elem).Name()
        }

        error.Report(error.NewError(error.RNT, span.Internal(), "Cannot Push() element of type '%s' into array of type '%s'!", typ, arr.Type.SubTypes[0].Name()))
        return nil
    }

    // append the new element
    arr.Elements = append(arr.Elements, elem)

    return nil
}

func Array_Pop(instance any, args []any) any {
    // make sure the instance isnt null
    if instance == nil {
        return 0
    }

    // get the last element of the array
    arr := instance.(*evalobjects.ArrayInstance)
    elem := arr.Elements[len(arr.Elements)-1]

    // remove the last element
    arr.Elements = arr.Elements[:len(arr.Elements)-1]

    return elem
}

// die(exitcode int)
func Die(args []any) any {
    // get the exit code
    code := args[0].(int32)
    os.Exit(int(code))

    return nil
}
