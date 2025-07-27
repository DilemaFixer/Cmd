package router

import (
	"fmt"
	"os"

	ctx "github.com/DilemaFixer/Cmd/context"
)

type Router struct {
	cache        map[string]RoutePoint
	points       map[string]RoutePoint
	errorHandler func(error, ctx.Context)
}

type RoutePoint interface {
	GetName() string
	Set(RoutePoint) error
	ProcessAndPush(ctx.Context, RoutingIterator) (RoutePoint, error)
}

func NewRouter() *Router {
	return &Router{
		cache:        make(map[string]RoutePoint),
		errorHandler: basicErrorHandler,
	}
}

func (r *Router) AddPoint(point RoutePoint) error {
	if point == nil {
		return fmt.Errorf("Router building errror: try add nil RoutePoint")
	}
	r.points[point.GetName()] = point
	return nil
}
