package router

import (
	"errors"
	"fmt"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/internal/utils"
	"strings"
)

/*
Router struct.
Maps: pattern -> method -> handlerFunc.
Eg.: /api/v1/user/:id -> GET -> GetUserHandler
*/
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
	pathStr := removeTrailingslash(r.URL.Path)
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
	possible := []route{}

	// check for possible matches
	for _, route := range router.routes {
		if len(route.pattern) != len(path) {
			continue
		}
		possible = append(possible, route)
	}

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

	found := false

	for _, route := range possible {
		found = true
		if route.method == method {
			variables := utils.GetPathVariables(path, route.pattern)
			handlerErr := route.handler(w, r, variables)
			if handlerErr != nil {
				logger.Log(logger.ERROR, "handler", handlerErr.Error())
				fmt.Println(handlerErr.Error())
			}
			return
		}
	}

	if found && method == "GET" {
		// send back webpage to browser
		MethodNotAllowed(w, r)
		return
	}

	// otherwise send back error
	Notfound(w, r)
	return
}

/*
removeTrailingslash removes trailing slash from path.
*/
func removeTrailingslash(path string) string {
	if path != "/" && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}
