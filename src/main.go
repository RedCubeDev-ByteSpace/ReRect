package main

import (
	"fmt"

	"bytespace.network/rerect/error"
	"bytespace.network/rerect/lexer"
	"bytespace.network/rerect/parser"
)

func main() {
    
    // Lexing
    // ------
    tokens := lexer.LexFile("./tests/print.rr")
    
    for _, v := range tokens {
        fmt.Printf("%s, %s, '%s'\n", v.Type, v.Position.Format(), v.Buffer)
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }

    // Parsing
    // -------
    members := parser.Parse(tokens)

    for _, v := range members {
        fmt.Printf("%s\n", v.Type())
    }

    // if there are errors -> output them and stop execution
    if error.HasErrors() {
        error.Output()
        return
    }
}
