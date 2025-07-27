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
	return fmt.Errorf("Can't set route point to EndPoint")
}

func (endPoint *EndPoint) ProcessAndPush(context ctx.Context, itr *RoutingIterator) (RoutePoint, error) {
	if err := endPoint.validateOptions(context); err != nil {
		return nil, err
	}

	return endPoint, endPoint.handler(context)
}

func (endPoint *EndPoint) validateOptions(context ctx.Context) error {
	if err := endPoint.validateGroups(context); err != nil {
		return err
	}

	if err := endPoint.validateGlobalOptions(context); err != nil {
		return err
	}

	return nil
}

func (endPoint *EndPoint) validateGroups(context ctx.Context) error {
	var solitudeGroupExist bool = false
	for _, group := range endPoint.groups.groups {
		if context.IsFlagExist(group.Triger) {
			if solitudeGroupExist {
				return fmt.Errorf("Routing error: Prev. group requires solitude, can't handling %s group", group.Triger)
			}

			solitudeGroupExist = group.RequiresSolitude
			if err := validateGroupOptions(group.Options, context); err != nil {
				return nil
			}
		}
	}
	return nil
}

func validateGroupOptions(options map[string]Option, context ctx.Context) error {
	for _, option := range options {
		var isExist bool
		if isExist = context.IsFlagExist(option.Name); isExist && option.Required {
			return fmt.Errorf("Routing error: Required %s flag not exist", option.Name)
		}

		if err := optionTypeValidation(option, context); err != nil {
			return err
		}
	}
	return nil
}

func (endPoint *EndPoint) validateGlobalOptions(context ctx.Context) error {
	for _, option := range endPoint.options {
		var isExist bool
		if isExist = context.IsFlagExist(option.Name); !isExist && option.Required {
			return fmt.Errorf("Routing error: Required %s flag not exist", option.Name)
		}

		if err := optionTypeValidation(option, context); err != nil {
			return err
		}
	}
	return nil
}

func optionTypeValidation(option Option, context ctx.Context) error {
	_type := option.Type
	switch _type {
	case Bool:
		if context.IsFlagHaveValue(option.Name) {
			return fmt.Errorf("Routing error: Option %s with type Bool have value, must look like --%s", option.Name, option.Name)
		}
	case String:
		if !context.IsFlagHaveValue(option.Name) {
			return fmt.Errorf("Routing error: Option %s with type String haven't value", option.Name)
		}
	case Int:
		if !context.IsFlagHaveValue(option.Name) {
			return fmt.Errorf("Routing error: Option %s with type Int haven't value", option.Name)
		}
		if _, error := context.GetValueAsInt(option.Name); error != nil {
			return fmt.Errorf("Routing error: Option %s with type Int have error \"%s\"", option.Name, error.Error())
		}
	case Float:
		if !context.IsFlagHaveValue(option.Name) {
			return fmt.Errorf("Routing error: Option %s with type Float haven't value", option.Name)
		}
		if _, error := context.GetValueAsFloat64(option.Name); error != nil {
			return fmt.Errorf("Routing error: Option %s with type Float have error \"%s\"", option.Name, error.Error())
		}
	default:
		return fmt.Errorf("Routing error: Undefine option type")
	}

	return nil
}
