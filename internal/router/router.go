package router

import (
	"errors"
	"fmt"
	"net/http"
	"pengoe/internal/logger"
	"pengoe/internal/utils"
	"reflect"
	"runtime"
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
ServeHTTP is mandatory
*/
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pathStr := r.URL.Path
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
	possible := router.routes

	for i, pathSegment := range path {
		newPossible := []route{}
		isMatch := false
		for _, route := range possible {
			patternSegment, rangeErr := utils.GetFromSlice(i, route.pattern)
			if rangeErr != nil {
				continue
			}

			if patternSegment == pathSegment {
				newPossible = append(newPossible, route)
				isMatch = true
				continue
			}
			if !isMatch && strings.HasPrefix(patternSegment, ":") {
				newPossible = append(newPossible, route)
			}
		}
		possible = newPossible
	}

	for _, route := range possible {
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

	if method == "GET" {
		// send back webpage to browser
		MethodNotAllowed(w, r)
		return
	}

	// otherwise send back error
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	return
}

/*
Url pattern matching.
Eg.: /api/v1/user/:id -> /api/v1/user/1/
*/
// func matches(pattern, path string) (bool, map[string]string) {
//
// 	// remove trailing slash
// 	pattern = removeTrailingslash(pattern)
// 	path = removeTrailingslash(path)
//
// 	patterns := strings.Split(pattern, "/")
// 	paths := strings.Split(path, "/")
//
// 	// check if the number of patterns and paths are the same
// 	if len(patterns) != len(paths) {
// 		return false, nil
// 	}
//
// 	pathVariables := make(map[string]string)
//
// 	// check if the patterns and paths match
// 	for i, path := range paths {
// 		pattern := patterns[i]
//
// 		if path == pattern {
// 			continue
// 		}
//
// 		if strings.HasPrefix(pattern, ":") {
// 			key := strings.TrimPrefix(pattern, ":")
// 			value := path
// 			pathVariables[key] = value
//
// 			continue
// 		}
//
// 		return false, nil
// 	}
//
// 	return true, pathVariables
// }

// remove trailing slash
func removeTrailingslash(path string) string {
	if path != "/" && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}

//list all routes with path, method and handler function name
func (router *Router) ListRoutes() {
	for i, routes := range router.routes {
		handlerFuncName := runtime.FuncForPC(reflect.ValueOf(routes.handler).Pointer()).Name()
		fmt.Printf("%d. %s %s -> %s\n", i+1, routes.method, routes.pattern, handlerFuncName)
	}
}
