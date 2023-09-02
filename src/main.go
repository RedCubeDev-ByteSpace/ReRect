package main

import (
	"fmt"

	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/error"
)

func main() {
    tokens := lexer.LexFile("./tests/print.rr")
    
    for _, v := range tokens {
        fmt.Printf("%s, %s, '%s'\n", v.Type, v.Position.Format(), v.Buffer)
    }


    if error.HasErrors() {
        error.Output()
    }
}
