package router

import (
	"errors"
	"fmt"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/internal/utils"
	"strings"
)

type Router struct {
	routes       []route
	staticPrefix string
	staticPath   string
}

type route struct {
	pattern []string
	method  string
	handler HandlerFunc
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string) error

/*
Utility function for creating a new router.
*/
func NewRouter() *Router {
	return &Router{
		routes: []route{},
	}
}

/*
SetStaticPath sets the static path for serving static files.
Accepts URL prefix and path to the static file directory.
*/
func (r *Router) SetStaticPath(prefix, path string) {
	r.staticPrefix = prefix
	r.staticPath = path
}

/*
Utility function for adding a new route to the router.
*/
func (r *Router) addRoute(method string, pattern []string, handler HandlerFunc) error {
	routeExists := false
	methodExists := false

	for _, route := range r.routes {
		if utils.SliceEqual(route.pattern, pattern) {
			routeExists = true
			if route.method == method {
				methodExists = true
				break
			}
		}
	}

	if routeExists && methodExists {
		return errors.New(fmt.Sprintf("Route %s %s already exists", method, pattern))
	}

	r.routes = append(r.routes, route{
		pattern,
		method,
		handler,
	})

	return nil
}

/*
Adds a new GET route to the router.
*/
func (r *Router) GET(s string, handler HandlerFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("GET", pattern, handler)
}

/*
Adds a new POST route to the router.
*/
func (r *Router) POST(s string, handler HandlerFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("POST", pattern, handler)
}

/*
Adds a new PUT route to the router.
*/
func (r *Router) PUT(s string, handler HandlerFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("PUT", pattern, handler)
}

/*
Adds a new DELETE route to the router.
*/
func (r *Router) DELETE(s string, handler HandlerFunc) {
	pattern := utils.GetPatternFromStr(s)
	r.addRoute("DELETE", pattern, handler)
}

/*
ServeHTTP is mandatory.
It searches for a matching route and calls the handler function.
*/
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathStr := removeTrailingSlash(r.URL.Path)
	method := r.Method

	// handle static files
	if router.staticPrefix != "" && strings.HasPrefix(pathStr, router.staticPrefix) {
		// w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		fs := http.FileServer(http.Dir(router.staticPath))
		staticHandler := http.StripPrefix(router.staticPrefix, fs)
		staticHandler.ServeHTTP(w, r)

		return
	}

	path := utils.GetPatternFromStr(pathStr)

	matchingRoutes, matchErr := matchRoutes(router.routes, path)
	if matchErr != nil {
		NotFound(w, r)
		return
	}

	route, methodErr := matchMethod(matchingRoutes, method)
	if methodErr != nil {
		if len(matchingRoutes) < 1 && method == "GET" {
			MethodNotAllowed(w, r)
			return
		}
		NotFound(w, r)
		return
	}

	variables := utils.GetPathVariables(route.pattern, path)

	handlerErr := route.handler(w, r, variables)
	if handlerErr != nil {
		logger.Log(logger.ERROR, "handler", handlerErr.Error())
	}
}

/*
getSameLengthRoutes returns routes with the same length as path.
*/
func getSameLengthRoutes(routes []route, path []string) []route {
	possible := []route{}

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
func matchRoutes(routes []route, path []string) ([]route, error) {
	possible := getSameLengthRoutes(routes, path)

	// get possible routes (should be only one)
	for i, pathSegment := range path {
		newPossible := []route{}
		// check for exact match
		for _, route := range possible {
			patternSegment := route.pattern[i]
			if pathSegment == patternSegment {
				newPossible = append(newPossible, route)
			}
		}
		// if no exact match, check for variable match
		if len(newPossible) == 0 {
			for _, route := range possible {
				patternSegment := route.pattern[i]
				if strings.HasPrefix(patternSegment, ":") {
					newPossible = append(newPossible, route)
				}
			}
		}
		possible = newPossible
	}

	if len(possible) == 0 {
		return []route{}, errors.New("No matching route found")
	}

	return possible, nil
}

/*
matchMethod returns the route that matches the method.
*/
func matchMethod(routes []route, method string) (route, error) {
	for _, route := range routes {
		if route.method == method {
			return route, nil
		}
	}

	return route{}, errors.New("No matching route found")
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
