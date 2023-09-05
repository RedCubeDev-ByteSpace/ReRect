package main

import (
	"bytespace.network/rerect/compctl"
	"bytespace.network/rerect/evaluator"
)

func main() {

    // Compile
    // -------
    prg := compctl.Compile("./tests/time.rr")

    // Evaluate 
    // --------
    evaluator.Evaluate(prg)
}
