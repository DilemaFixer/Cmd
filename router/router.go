package router

import (
	"fmt"

	ctx "github.com/DilemaFixer/Cmd/context"
)

type Router struct {
	cache  map[string]RoutePoint
	points []RoutePoint
}

type RoutePoint interface {
	GetName() string
	Set(RoutePoint) error
	ProcessAndPush(ctx.Context, RoutingIterator) error
}

func NewRouter() *Router {
	return &Router{
		cache: make(map[string]RoutePoint),
	}
}

func (r *Router) AddPoint(point RoutePoint) error {
	if point == nil {
		return fmt.Errorf("Router building errror: try add nil RoutePoint")
	}
	r.points = append(r.points, point)
	return nil
}
