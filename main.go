package main

import (
	"fmt"
	p "github.com/DilemaFixer/Cmd/parser"
)

func main() {
	input := "command subcommand --bool_flag --value-flag=12"
	result, err := p.ParseInput(input)
	if err != nil {
		fmt.Println(err)
	}else {
		fmt.Println(result)
	}
}
