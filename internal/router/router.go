package router

import (
	"errors"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/internal/utils"
	"strings"
)

type Router struct {
	Routes       []*Route
	StaticPrefix string
	StaticPath   string
}

type Route struct {
	Pattern     []string
	Method      string
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

/*
Utility function for creating a new router.
*/
func NewRouter() *Router {
	return &Router{
		Routes:       []*Route{},
		StaticPrefix: "",
		StaticPath:   "",
	}
}

/*
SetStaticPath sets the static path for serving static files.
Accepts URL prefix and path to the static file directory.
*/
func (r *Router) SetStaticPath(prefix, path string) {
	r.StaticPrefix = prefix
	r.StaticPath = path
}

/*
Utility function for adding a new route to the router.
*/
func (r *Router) AddRoute(method string, pattern []string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	for _, route := range r.Routes {
		if utils.SliceEqual(route.Pattern, pattern) && route.Method == method {
			return
		}
	}

	newRoute := &Route{
		pattern,
		method,
		handler,
		middlewares,
	}

	r.Routes = append(r.Routes, newRoute)
}

/*
Adds a new GET route to the router.
*/
func (r *Router) GET(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.AddRoute("GET", pattern, handler, middlewares...)
}

/*
Adds a new POST route to the router.
*/
func (r *Router) POST(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.AddRoute("POST", pattern, handler, middlewares...)
}

/*
Adds a new PATCH route to the router.
*/
func (r *Router) PATCH(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.AddRoute("PATCH", pattern, handler, middlewares...)
}

/*
Adds a new DELETE route to the router.
*/
func (r *Router) DELETE(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.AddRoute("DELETE", pattern, handler, middlewares...)
}

/*
ServeHTTP is mandatory.
It searches for a matching route and calls the handler function.
*/
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathStr := utils.RemoveTrailingSlash(req.URL.Path)
	method := req.Method

	// handle static files
	if r.StaticPrefix != "" && strings.HasPrefix(pathStr, r.StaticPrefix) {
		// w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fs := http.FileServer(http.Dir(r.StaticPath))
		staticHandler := http.StripPrefix(r.StaticPrefix, fs)
		staticHandler.ServeHTTP(w, req)

		return
	}

	// /account/1 -> [account :id]
	path := utils.GetPatternFromStr(pathStr)

	// same length routes
	sameLengthRoutes := GetSameLengthRoutes(r.Routes, path)

	// [account :id] -> [account 1]
	matchingRoutes, matchErr := MatchRoutes(sameLengthRoutes, path)
	if matchErr != nil {
		NotFound(w, req, nil)
		return
	}

	// GET -> GET
	route, methodErr := MatchMethod(matchingRoutes, method)
	if methodErr != nil {
		if len(matchingRoutes) < 1 && method == "GET" {
			MethodNotAllowed(w, req, nil)
		} else {
			NotFound(w, req, nil)
		}
		return
	}

	// [account 1], [account :id] -> {id: 1}
	variables := utils.GetPathVariables(route.Pattern, path)

	handler := route.Handler

	// apply middlewares backwards
	for i := len(route.Middlewares) - 1; i >= 0; i-- {
		handler = route.Middlewares[i](handler)
	}

	// call handler
	handlerErr := handler(w, req, variables)
	if handlerErr != nil {
		logger.Error(handlerErr.Error())
	}
}

/*
GetSameLengthRoutes returns routes with the same length as path.
*/
func GetSameLengthRoutes(routes []*Route, path []string) []*Route {
	possible := []*Route{}

	for _, route := range routes {
		if len(route.Pattern) != len(path) {
			continue
		}
		possible = append(possible, route)
	}

	return possible
}

/*
MatchRoutes returns the route that matches the path.
Works only for same length routes.
You must filter routes by getSameLengthRoutes first.
*/
func MatchRoutes(routes []*Route, path []string) ([]*Route, error) {
	result := routes

	// get possible routes (should be only one)
	for i, pathSegment := range path {
		newPossible := []*Route{}
		// check for exact match
		for _, route := range result {
			patternSegment := route.Pattern[i]
			if pathSegment == patternSegment {
				newPossible = append(newPossible, route)
			}
		}
		// if no exact match, check for variable match
		if len(newPossible) == 0 {
			for _, route := range result {
				patternSegment := route.Pattern[i]
				if strings.HasPrefix(patternSegment, ":") {
					newPossible = append(newPossible, route)
				}
			}
		}
		result = newPossible
	}

	if len(result) == 0 {
		return nil, errors.New("No matching route found")
	}

	return result, nil
}

/*
MatchMethod returns the route that matches the method.
*/
func MatchMethod(routes []*Route, method string) (*Route, error) {
	for _, route := range routes {
		if route.Method == method {
			return route, nil
		}
	}

	return nil, errors.New("No matching route found")
}
