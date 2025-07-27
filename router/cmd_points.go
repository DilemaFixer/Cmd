package router

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
)

type CmdPoint struct {
	points map[string]RoutePoint
	name   string
}

func NewCmdPoint(name string) *CmdPoint {
	return &CmdPoint{
		points: make(map[string]RoutePoint),
		name:   name,
	}
}

func (cmd *CmdPoint) GetName() string {
	return cmd.name
}

func (cmd *CmdPoint) Set(point RoutePoint) error {
	name := point.GetName()
	_, exist := cmd.points[name]

	if exist {
		return fmt.Errorf("Router building error: Can't add new route point with name %s to %s , it exist", name, cmd.name)
	}

	cmd.points[name] = point
	return nil
}

func (cmd *CmdPoint) ProcessAndPush(context ctx.Context, itr *RoutingIterator) (RoutePoint, error) {
	//TODO: how i can paste routing iterator more beauti
	next, exist := cmd.points[itr.Get()]
	if !exist {
		return nil, fmt.Errorf("Routing error: Point with name %s not found ", itr.Get())
	}
	itr.Next()
	return next.ProcessAndPush(context, itr)
}

func (c *CmdPoint) AddSubCommand(name string, point RoutePoint) {
	c.points[name] = point
}

func (c *CmdPoint) GetSubCommand(name string) (RoutePoint, bool) {
	point, exists := c.points[name]
	return point, exists
}

func (c *CmdPoint) GetAllSubCommands() map[string]RoutePoint {
	return c.points
}
