package context

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	prs "github.com/DilemaFixer/Cmd/parser"
)

type Context struct {
	command     string
	subcommands map[string]struct{}
	flags       map[string]string
}

func NewContext(input *prs.ParsedInput) *Context {
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
	if strings.HasPrefix(name, "--") {
		name = name[2:]
	}
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

func (ctx *Context) GetValueAsBool(name string) bool {
	_, exists := ctx.flags[name]
	if !exists {
		return false
	}

	return true
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

func (ctx *Context) GetCommand() string {
	return ctx.command
}

func (ctx *Context) IsCommandEqual(target string) bool {
	return ctx.command == target
}

func (ctx *Context) IsSubcommandExist(target string) bool {
	_, exist := ctx.subcommands[target]
	return exist
}

func (ctx *Context) GetSubcommandsAsArr() []string {
	subcommandsArr := make([]string, 0)

	if len(ctx.subcommands) == 0 {
		return subcommandsArr
	}

	for subcommand, _ := range ctx.subcommands {
		subcommandsArr = append(subcommandsArr, subcommand)
	}

	return subcommandsArr
}

func (ctx *Context) GetFlagsAsMap() map[string]string {
	flagsMap := make(map[string]string)

	for flag, value := range ctx.flags {
		flagsMap[flag] = value
	}

	return flagsMap
}

func (ctx *Context) GetFlagsAsArr() []string {
	flagsArr := make([]string, 0)

	for flag, value := range ctx.flags {
		flagsArr = append(flagsArr, fmt.Sprintf("%s=%s", flag, value))
	}

	return flagsArr
}

func (ctx *Context) GetFlagsKeysAsArr() []string {
	flagsKeysArr := make([]string, 0)

	for flag, _ := range ctx.flags {
		flagsKeysArr = append(flagsKeysArr, flag)
	}

	return flagsKeysArr
}

func (ctx *Context) GetFlagsValuesAsArr() []string {
	flagsValuesArr := make([]string, 0)

	for _, value := range ctx.flags {
		flagsValuesArr = append(flagsValuesArr, value)
	}

	return flagsValuesArr
}
