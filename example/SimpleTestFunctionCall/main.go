package main

import (
	"fmt"
	ctx "github.com/DilemaFixer/Cmd/context"
	p "github.com/DilemaFixer/Cmd/parser"
	rtr "github.com/DilemaFixer/Cmd/router"
)

func main() {
	input := "command subcommand sub endpoint --bool_flag --value-flag=12"
	parsedInput, err := p.ParseInput(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(parsedInput)
	}

	context := ctx.NewContext(parsedInput)
	itr := rtr.NewRoutingIterator(context)
	router := rtr.NewRouter()

	router.NewCmd("command").
		NewSub("subcommand").
		NewSub("sub").
		Endpoint("endpoint").
		RequiredBool("bool_flag").
		RequiredInt("value-flag").
		Handler(test).
		Build().
		Build().
		Build().
		Register()

	router.Route(*context, itr)
}

func test(context ctx.Context) error {
	println("test func is calling")
	return nil
}
