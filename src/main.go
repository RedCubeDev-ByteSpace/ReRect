package main

import (
	"fmt"
	"os"

	"bytespace.network/rerect/compctl"
	"bytespace.network/rerect/evaluator"
	"bytespace.network/rerect/error"
)

func main() {

    if len(os.Args) < 2 {
        fmt.Println("At least one source file required!")
        return
    }


    // Compile
    // -------
    prg := compctl.Compile(os.Args[1:])

    if !prg.Ok {
        return
    }

    // Evaluate 
    // --------
    evaluator.Evaluate(prg)

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }
}
