// Go packages - sys.go
// --------------------------------------------------------
// ah yes, Rects sys library, its good to be back
// --------------------------------------------------------
package gopackages

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

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

    /* sys::Write() */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Write", compunit.GlobalDataTypeRegister["void"]  , []*symbols.ParameterSymbol{symbols.NewParameterSymbol("msg", 0, compunit.GlobalDataTypeRegister["string"])}, Write))
    /* sys::Input() */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Input", compunit.GlobalDataTypeRegister["string"], []*symbols.ParameterSymbol{}, Input))
    /* sys::Clear() */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Clear", compunit.GlobalDataTypeRegister["void"]  , []*symbols.ParameterSymbol{}, Clear))
    /* sys::Sleep() */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Sleep", compunit.GlobalDataTypeRegister["void"]  , []*symbols.ParameterSymbol{symbols.NewParameterSymbol("mills", 0, compunit.GlobalDataTypeRegister["long"])}, Sleep))
    /* sys::Now()   */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Now"  , compunit.GlobalDataTypeRegister["long"]  , []*symbols.ParameterSymbol{}, Now))
    /* sys::Char()  */ registerFunction("sys", symbols.NewVMFunctionSymbol(sys, "Char" , compunit.GlobalDataTypeRegister["string"], []*symbols.ParameterSymbol{symbols.NewParameterSymbol("ascii", 0, compunit.GlobalDataTypeRegister["int"])}, Char))
}

// sys::Print(msg string)
func Print(args []any) any {
    // upack args
    msg := args[0].(string)

    // do the thing
    fmt.Println(msg)

    // ok we don
    return nil
}

// sys::Write(msg string)
func Write(args []any) any {
    // upack args
    msg := args[0].(string)

    // do the thing
    fmt.Print(msg)

    // ok we don
    return nil
}

// sys::Input() string
func Input(args []any) any {
    // do the thing
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')

    // remove \n and \r
    input = strings.Replace(input, "\n", "", -1)
    input = strings.Replace(input, "\r", "", -1)

    // ok we don
    return input
}

// sys::Clear()
func Clear(args []any) any {
    // do the thing
    fmt.Print("\033[H\033[2J")

    // ok we don
    return nil
}

// sys::Sleep(mills long) 
func Sleep(args []any) any {
    // unpack args
    mills := args[0].(int64)

    // do the thing
    time.Sleep(time.Duration(mills) * time.Millisecond)

    // ok we don
    return nil
}

// sys::Now() long 
func Now(args []any) any {
    // do the thing
    return time.Now().UnixMilli()
}

// sys::Char(ascii int) string 
func Char(args []any) any {
    // unpack args
    ascii := args[0].(int32)

    // do the thing
    return string(rune(ascii))
}

// sys::die(exitcode int) 

