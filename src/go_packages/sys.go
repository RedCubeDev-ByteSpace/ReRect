// Go packages - sys.go
// --------------------------------------------------------
// ah yes, Rects sys library, its good to be back
// --------------------------------------------------------
package gopackages

import (
	"fmt"

	"bytespace.network/rerect/compunit"
	"bytespace.network/rerect/symbols"
)

func LoadSys() {
    sys := registerPackage("sys")
    registerFunction("sys", symbols.NewVMFunctionSymbol(
        // The package this function belongs to
        sys,

        // Function name
        "Print",

        // Return type
        compunit.GlobalDataTypeRegister["void"],

        // Function parameters
        []*symbols.ParameterSymbol {
            symbols.NewParameterSymbol("msg", 0, compunit.GlobalDataTypeRegister["string"]),
        },

        // Pointer to function
        Print,
    ))
}

// sys::Print(string)
func Print(args []any) any {
    fmt.Println(args[0].(string))
    return nil
}
