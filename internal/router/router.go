package router

import (
	"errors"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/internal/utils"
	"strings"
)

type router struct {
	routes       []*route
	staticPrefix string
	staticPath   string
}

type route struct {
	pattern []string
	method  string
	handler HandlerFunc
	middlewares []MiddlewareFunc
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

/*
Utility function for creating a new router.
*/
func NewRouter() *router {
	return &router{
		routes: []*route{},
		staticPrefix: "",
		staticPath: "",
	}
}

/*
SetStaticPath sets the static path for serving static files.
Accepts URL prefix and path to the static file directory.
*/
func (r *router) SetStaticPath(prefix, path string) {
	r.staticPrefix = prefix
	r.staticPath = path
}

/*
Utility function for adding a new route to the router.
*/
func (r *router) addRoute(method string, pattern []string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	for _, route := range r.routes {
		if utils.SliceEqual(route.pattern, pattern) && route.method == method {
			return
		}
	}

	newRoute := &route{
		pattern,
		method,
		handler,
		middlewares,
	}

	r.routes = append(r.routes, newRoute)
}

/*
Adds a new GET route to the router.
*/
func (r *router) GET(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("GET", pattern, handler, middlewares...)
}

/*
Adds a new POST route to the router.
*/
func (r *router) POST(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("POST", pattern, handler, middlewares...)
}

/*
Adds a new PUT route to the router.
*/
func (r *router) PUT(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("PUT", pattern, handler, middlewares...)
}

/*
Adds a new DELETE route to the router.
*/
func (r *router) DELETE(s string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("DELETE", pattern, handler, middlewares...)
}

/*
ServeHTTP is mandatory.
It searches for a matching route and calls the handler function.
*/
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pathStr := removeTrailingSlash(req.URL.Path)
	method := req.Method

	// handle static files
	if r.staticPrefix != "" && strings.HasPrefix(pathStr, r.staticPrefix) {
		// w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fs := http.FileServer(http.Dir(r.staticPath))
		staticHandler := http.StripPrefix(r.staticPrefix, fs)
		staticHandler.ServeHTTP(w, req)

		return
	}

	// /account/1 -> [account :id]
	path := utils.GetPatternFromStr(pathStr)

	// same length routes
	sameLengthRoutes := getSameLengthRoutes(r.routes, path)

	// [account :id] -> [account 1]
	matchingRoutes, matchErr := matchRoutes(sameLengthRoutes, path)
	if matchErr != nil {
		NotFound(w, req, nil)
		return
	}

	// GET -> GET
	route, methodErr := matchMethod(matchingRoutes, method)
	if methodErr != nil {
		if len(matchingRoutes) < 1 && method == "GET" {
			MethodNotAllowed(w, req, nil)
		} else {
			NotFound(w, req, nil)
		}
		return
	}

	// [account 1], [account :id] -> {id: 1}
	variables := utils.GetPathVariables(route.pattern, path)

	handler := route.handler

	// apply middlewares
	for _, middleware := range route.middlewares {
		handler = middleware(handler)
	}

	// call handler
	handlerErr := handler(w, req, variables)
	if handlerErr != nil {
		logger.Log(logger.ERROR, "handler", handlerErr.Error())
	}
}







/*
getSameLengthRoutes returns routes with the same length as path.
*/
func getSameLengthRoutes(routes []*route, path []string) []*route {
	possible := []*route{}

	for _, route := range routes {
		if len(route.pattern) != len(path) {
			continue
		}
		possible = append(possible, route)
	}

	return possible
}

/*
matchRoutes returns the route that matches the path.
Works only for same length routes.
You must filter routes by getSameLengthRoutes first.
*/
func matchRoutes(routes []*route, path []string) ([]*route, error) {
	result := routes

	// get possible routes (should be only one)
	for i, pathSegment := range path {
		newPossible := []*route{}
		// check for exact match
		for _, route := range result {
			patternSegment := route.pattern[i]
			if pathSegment == patternSegment {
				newPossible = append(newPossible, route)
			}
		}
		// if no exact match, check for variable match
		if len(newPossible) == 0 {
			for _, route := range result {
				patternSegment := route.pattern[i]
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
matchMethod returns the route that matches the method.
*/
func matchMethod(routes []*route, method string) (*route, error) {
	for _, route := range routes {
		if route.method == method {
			return route, nil
		}
	}

	return nil, errors.New("No matching route found")
}

/*
removeTrailingSlash removes trailing slash from path.
*/
func removeTrailingSlash(path string) string {
	if path != "/" && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}
