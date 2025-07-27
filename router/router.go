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
func basicErrorHandler(err error, ctx ctx.Context) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
	os.Exit(1)
}

func (r *Router) CustomErrorHandler(errorHandler func(error, ctx.Context)) {
	if errorHandler == nil {
		return
	}
	r.errorHandler = errorHandler
}

func (r *Router) Route(context ctx.Context, itr RoutingIterator) {
	point, exist := r.points[itr.RouteToString()]

	if exist {
		if _, err := point.ProcessAndPush(context, itr); err != nil {
			r.errorHandler(err, context)
		}
		return
	}

	point, exist = r.points[itr.Get()]

	if !exist {
		r.errorHandler(fmt.Errorf("Router routing error: try route to non-existent point"), context)
		return
	}

	point, err := point.ProcessAndPush(context, itr)
	if err != nil {
		r.errorHandler(err, context)
	}
	r.cache[point.GetName()] = point
}
