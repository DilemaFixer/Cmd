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

func (endPoint *EndPoint) ProcessAndPush(context ctx.Context, itr RoutingIterator) (RoutePoint, error) {
	if err := validateOptions(endPoint, context); err != nil {
		return nil, err
	}

	return endPoint, endPoint.handler(context)
}

func validateOptions(endPoint *EndPoint, context ctx.Context) error {

	var throwError bool = false // change name to more readable

	//TODO: refactor this part , move to new function
	for _, value := range endPoint.groups.groups {
		if context.IsFlagExist(value.Triger) {
			if throwError {
				return nil // DOTO: Throw error if flag must be lonely
			}

			throwError = value.RequiresSolitude
			for _, option := range value.Options {
				var isExist bool
				if isExist = context.IsFlagExist(option.Name); isExist && option.Required {
					return nil //TODO: throw error that requared flag is exist
				}

				if option.Type != Bool && isExist {
					if !context.IsFlagHaveValue(option.Name) {
						return nil // TODO throw error that invalid type or value exist
					}
				}
			}
		}
	}

	for _, option := range endPoint.options {
		var isExist bool
		if isExist = context.IsFlagExist(option.Name); !isExist && option.Required {
			return nil // DOTO: throwing error
		}

		if isExist && option.Type != Bool {
			if !context.IsFlagHaveValue(option.Name) {
				return nil //TODO: throwing err
			}
		}
	}
	return nil
}

func optionTypeValidation(option Option, context ctx.Context, _type OptionType) error {
	switch _type {
	case Bool:
		if context.IsFlagHaveValue(option.Name) {
			return nil //TODO: thr err
		}
	case String:
		if !context.IsFlagHaveValue(option.Name) {
			return nil //TODO: thr err
		}
	case Int:
		if _, error := context.GetValueAsInt(option.Name); error != nil {
			return nil //TODO: format and return err
		}
	case Float:
		if _, error := context.GetValueAsFloat64(option.Name); error != nil {
			return nil //TODO: format and return err
		}
	default:
		return nil //return undefind
	}

	return nil
}
