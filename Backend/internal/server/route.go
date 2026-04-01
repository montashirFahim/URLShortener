package server

import (
	"net/http"
)

// Route represents a Route
type Route struct {
	Pattern string
	Handler http.Handler
	Method  string
}

// Routes is a list of Route
type Routes []Route

// Add adds a route to Routes
func (rs *Routes) Add(route Route) {
	*rs = append(*rs, route)
}

func (rs Routes) ApplyRouteModifire(modifiers ...RouteModifier) {
	for i := range rs {
		for _, modifier := range modifiers {
			modifier(&rs[i])
		}
	}
}
