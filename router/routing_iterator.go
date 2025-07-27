package router

import (
	ctx "github.com/DilemaFixer/Cmd/context"
)

type RoutingIterator struct {
	i    int
	maxI int
	rout []string
}

func NewRoutingIterator(context *ctx.Context) RoutingIterator {
	route, count := concatCommandAndSubcommands(context)
	return RoutingIterator{
		i:    0,
		maxI: count - 1,
		rout: route,
	}
}

func concatCommandAndSubcommands(context *ctx.Context) ([]string, int) {
	command := context.GetCommand()
	subcommands := context.GetSubcommandsAsArr()
	count := len(subcommands) + 1
	result := make([]string, count)
	result = append(result, command)
	for _, value := range subcommands {
		result = append(result, value)
	}
	return result, count
}

func (itr RoutingIterator) Get() string {
	return itr.rout[itr.i]
}

func (itr RoutingIterator) Next() bool {
	if itr.maxI == itr.i {
		return false
	}
	itr.i++
	return true
}

func (itr RoutingIterator) CheckOnTarget(point RoutePoint) bool {
	if itr.Get() == point.GetName() {
		return true
	}
	return false
}
