package context

import (
	"errors"
	"strconv"

	prs "github.com/DilemaFixer/Cmd/parser"
)

type Context struct {
	command     string
	subcommands map[string]struct{}
	flags       map[string]string
}

func NewContext(input prs.ParserInput) *Context {
	ctx := &Context{
		command:     input.Command,
		subcommands: make(map[string]struct{}),
		flags:       make(map[string]string),
	}

	for _, subcommand := range input.Subcommands {
		ctx.subcommands[subcommand] = struct{}{}
	}

	for _, flag := range input.InputFlags {
		ctx.flags[flag.Name] = flag.Value
	}

	return ctx
}

func (ctx *Context) IsFlagExist(name string) bool {
	_, exists := ctx.flags[name]
	return exists
}

func (ctx *Context) IsFlagHaveValue(name string) bool {
	value, exists := ctx.flags[name]
	return exists && value != ""
}

func (ctx *Context) GetValueAsInt32(name string) (int32, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return 0, errors.New("flag not found")
	}
	if value == "" {
		return 0, errors.New("flag has empty value")
	}

	parsed, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

func (ctx *Context) GetValueAsInt64(name string) (int64, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return 0, errors.New("flag not found")
	}
	if value == "" {
		return 0, errors.New("flag has empty value")
	}

	return strconv.ParseInt(value, 10, 64)
}

func (ctx *Context) GetValueAsInt(name string) (int, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return 0, errors.New("flag not found")
	}
	if value == "" {
		return 0, errors.New("flag has empty value")
	}

	return strconv.Atoi(value)
}

func (ctx *Context) GetValueAsFloat32(name string) (float32, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return 0, errors.New("flag not found")
	}
	if value == "" {
		return 0, errors.New("flag has empty value")
	}

	parsed, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return 0, err
	}
	return float32(parsed), nil
}

func (ctx *Context) GetValueAsFloat64(name string) (float64, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return 0, errors.New("flag not found")
	}
	if value == "" {
		return 0, errors.New("flag has empty value")
	}

	return strconv.ParseFloat(value, 64)
}

func (ctx *Context) GetValueAsBool(name string) (bool, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return false, errors.New("flag not found")
	}
	if value == "" {
		return false, errors.New("flag has empty value")
	}

	return strconv.ParseBool(value)
}

func (ctx *Context) GetValueAsString(name string) (string, error) {
	value, exists := ctx.flags[name]
	if !exists {
		return "", errors.New("flag not found")
	}
	return value, nil
}

func (ctx *Context) GetValueOrDefault(name, defaultValue string) string {
	if value, exists := ctx.flags[name]; exists && value != "" {
		return value
	}
	return defaultValue
}
