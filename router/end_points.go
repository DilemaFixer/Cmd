package router

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
)

type OptionType int

const (
	Bool OptionType = iota
	String
	Int
	Float
)

type EndPoint struct {
	name        string
	handler     func(ctx.Context) error
	options     map[string]Option
	groups      OptionsGroups
	description string
}

type OptionsGroups struct {
	CanBeIgnored bool
	groups       map[string]OptionsGroup
}

type OptionsGroup struct {
	Triger           string
	RequiresSolitude bool
	Options          map[string]Option
}

type Option struct {
	Name     string
	Type     OptionType
	Required bool
}

func NewEndPoint(name string, handler func(ctx.Context) error) *EndPoint {
	return &EndPoint{
		name:        name,
		handler:     handler,
		options:     make(map[string]Option),
		groups:      NewOptionsGroups(),
		description: "",
	}
}

func NewOptionsGroups() OptionsGroups {
	return OptionsGroups{
		CanBeIgnored: false,
		groups:       make(map[string]OptionsGroup),
	}
}

func NewOptionsGroup(trigger string, requiresSolitude bool) OptionsGroup {
	return OptionsGroup{
		Triger:           trigger,
		RequiresSolitude: requiresSolitude,
		Options:          make(map[string]Option),
	}
}

func NewOption(name string, optionType OptionType, required bool) Option {
	return Option{
		Name:     name,
		Type:     optionType,
		Required: required,
	}
}

func (endPoint *EndPoint) GetName() string {
	return endPoint.name
}

func (endPoint *EndPoint) Set(routePoint RoutePoint) error {
	return fmt.Errorf("not implemented")
}

func (endPoint *EndPoint) ProcessAndPush(ctx.Context, RoutingIterator) error {
	return fmt.Errorf("not implemented")
}
