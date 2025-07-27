package router

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
)

type CmdWrapper struct {
	router *Router
	cmd    *CmdPoint
	parent *CmdWrapper
}

type EndPointWrapper struct {
	router   *Router
	endpoint *EndPoint
	parent   *CmdWrapper
}

type EndPointGroupWrapper struct {
	endpointWrapper *EndPointWrapper
	groupName       string
}

func (r *Router) NewCmd(name string) *CmdWrapper {
	return &CmdWrapper{
		router: r,
		cmd:    NewCmdPoint(name),
	}
}

func (cmd *CmdWrapper) NewSub(name string) *CmdWrapper {
	subCmd := NewCmdPoint(name)
	cmd.cmd.AddSubCommand(name, subCmd)
	return &CmdWrapper{
		router: cmd.router,
		cmd:    subCmd,
		parent: cmd,
	}
}

func (cmd *CmdWrapper) Endpoint(name string) *EndPointWrapper {
	endpoint := NewEndPoint(name, nil)
	cmd.cmd.AddSubCommand(name, endpoint)

	return &EndPointWrapper{
		router:   cmd.router,
		endpoint: endpoint,
		parent:   cmd,
	}
}

func (r *Router) Endpoint(name string) *EndPointWrapper {
	return &EndPointWrapper{
		router:   r,
		endpoint: NewEndPoint(name, nil),
	}
}

func (cmd *CmdWrapper) Build() *CmdWrapper {
	if cmd.parent == nil {
		cmd.router.AddPoint(cmd.cmd)
		return cmd
	}
	return cmd.parent
}

func (cmd *CmdWrapper) Register() {
	fmt.Println(cmd.cmd)
	cmd.router.AddPoint(cmd.cmd)
}

func (w *EndPointWrapper) Description(desc string) *EndPointWrapper {
	w.endpoint.description = desc
	return w
}

func (w *EndPointWrapper) Handler(handler func(ctx.Context) error) *EndPointWrapper {
	w.endpoint.handler = handler
	return w
}

func (w *EndPointWrapper) Option(name string, optType OptionType, required bool) *EndPointWrapper {
	w.endpoint.options[name] = NewOption(name, optType, required)
	return w
}

func (w *EndPointWrapper) StringOption(name string) *EndPointWrapper {
	return w.Option(name, String, false)
}

func (w *EndPointWrapper) RequiredString(name string) *EndPointWrapper {
	return w.Option(name, String, true)
}

func (w *EndPointWrapper) IntOption(name string) *EndPointWrapper {
	return w.Option(name, Int, false)
}

func (w *EndPointWrapper) RequiredInt(name string) *EndPointWrapper {
	return w.Option(name, Int, true)
}

func (w *EndPointWrapper) BoolOption(name string) *EndPointWrapper {
	return w.Option(name, Bool, false)
}

func (w *EndPointWrapper) RequiredBool(name string) *EndPointWrapper {
	return w.Option(name, Bool, true)
}

func (w *EndPointWrapper) FloatOption(name string) *EndPointWrapper {
	return w.Option(name, Float, false)
}

func (w *EndPointWrapper) RequiredFloat(name string) *EndPointWrapper {
	return w.Option(name, Float, true)
}

func (w *EndPointWrapper) Group(name, trigger string) *EndPointGroupWrapper {
	group := NewOptionsGroup(trigger, false)
	w.endpoint.groups.groups[name] = group

	return &EndPointGroupWrapper{
		endpointWrapper: w,
		groupName:       name,
	}
}

func (w *EndPointWrapper) ExclusiveGroup(name, trigger string) *EndPointGroupWrapper {
	group := NewOptionsGroup(trigger, true)
	w.endpoint.groups.groups[name] = group

	return &EndPointGroupWrapper{
		endpointWrapper: w,
		groupName:       name,
	}
}

func (w *EndPointWrapper) SetGroupsCanBeIgnored(canBeIgnored bool) *EndPointWrapper {
	w.endpoint.groups.CanBeIgnored = canBeIgnored
	return w
}

func (w *EndPointWrapper) Build() *CmdWrapper {
	if w.endpoint.handler == nil {
		panic(fmt.Sprintf("Error router building: endpoint \"%s\" handler is not set", w.endpoint.name))
	}

	if w.parent == nil {
		w.router.AddPoint(w.endpoint)
		return nil
	} else {
		w.parent.cmd.AddSubCommand(w.endpoint.name, w.endpoint)
	}
	return w.parent
}

func (w *EndPointWrapper) Register() {
	if w.parent == nil {
		w.router.AddPoint(w.endpoint)
	}
}

func (w *EndPointGroupWrapper) GroupOption(name string, optType OptionType, required bool) *EndPointGroupWrapper {
	group := w.endpointWrapper.endpoint.groups.groups[w.groupName]
	group.Options[name] = NewOption(name, optType, required)
	w.endpointWrapper.endpoint.groups.groups[w.groupName] = group
	return w
}

func (w *EndPointGroupWrapper) StringOption(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, String, false)
}

func (w *EndPointGroupWrapper) RequiredString(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, String, true)
}

func (w *EndPointGroupWrapper) IntOption(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Int, false)
}

func (w *EndPointGroupWrapper) RequiredInt(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Int, true)
}

func (w *EndPointGroupWrapper) BoolOption(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Bool, false)
}

func (w *EndPointGroupWrapper) RequiredBool(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Bool, true)
}

func (w *EndPointGroupWrapper) FloatOption(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Float, false)
}

func (w *EndPointGroupWrapper) RequiredFloat(name string) *EndPointGroupWrapper {
	return w.GroupOption(name, Float, true)
}

func (w *EndPointGroupWrapper) EndGroup() *EndPointWrapper {
	return w.endpointWrapper
}
