package router

import (
	"net/http"
	"pengoe/internal/utils"
	"testing"
)

func TestAddRoute(t *testing.T) {
	r := NewRouter()

	var handler HandlerFunc = func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		return nil
	}

	r.addRoute("GET", []string{"test", ":id", "helo"}, handler)

	if r.routes[0].method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", r.routes[0].method)
	}

	if utils.SliceEqual(r.routes[0].pattern, []string{"test", ":id", "helo"}) == false {
		t.Errorf("Expected pattern '[test :id helo]', got '%s'", r.routes[0].pattern)
	}

	if r.routes[0].handler == nil {
		t.Errorf("Expected handler address to not be nil")
	}
}

func TestRemoveTrailingSlash(t *testing.T) {
	path := "/test/path/"
	expected := "/test/path"

	result := removeTrailingSlash(path)

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGetSameLengthRoutes(t *testing.T) {
	routes := []*route{
		{
			pattern: []string{"one", "two"},
			method:  "GET",
			handler: nil,
		},
		{
			pattern: []string{"one", "two", "three"},
			method:  "GET",
			handler: nil,
		},
		{
			pattern: []string{"one", "two", "other"},
			method:  "GET",
			handler: nil,
		},
		{
			pattern: []string{"one", "two", "three", "four"},
			method:  "GET",
			handler: nil,
		},
	}

	path := []string{"one", "two", "three"}

	result := getSameLengthRoutes(routes, path)

	if len(result) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(result))
	}

	for _, route := range result {
		if len(route.pattern) != len(path) {
			t.Errorf("Expected route pattern length to be %d, got %d", len(path), len(route.pattern))
		}
	}
}

func TestMatchRoutes(t *testing.T) {
	routes := []*route{
		{
			pattern: []string{"one", "two"},
			method:  "GET",
			handler: nil,
		},
		{
			pattern: []string{"one", ":var"},
			method:  "GET",
			handler: nil,
		},
	}

	path1 := []string{"one", "two"}
	result1, err1 := matchRoutes(routes, path1)
	if err1 != nil {
		t.Errorf("Expected no error, got '%s'", err1)
	}
	if len(result1) != 1 {
		t.Errorf("Expected 1 route, got %d", len(result1))
	}
	if !utils.SliceEqual(result1[0].pattern, []string{"one", "two"}) {
		t.Errorf("Expected route pattern to be '[one two]', got '%s'", result1[0].pattern)
	}

	path2 := []string{"one", "three"}
	result2, err2 := matchRoutes(routes, path2)
	if err2 != nil {
		t.Errorf("Expected no error, got '%s'", err2)
	}
	if len(result2) != 1 {
		t.Errorf("Expected 1 route, got %d", len(result2))
	}
	if !utils.SliceEqual(result2[0].pattern, []string{"one", ":var"}) {
		t.Errorf("Expected route pattern to be '[one :var]', got '%s'", result2[0].pattern)
	}

	path3 := []string{"two", "three"}
	result3, err3 := matchRoutes(routes, path3)
	if err3 == nil {
		t.Errorf("Expected error, got nil")
	}
	if len(result3) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(result3))
	}
}

func TestMatchMethod(t *testing.T) {
	routes := []*route{
		{
			pattern: []string{"one", "two"},
			method:  "GET",
			handler: nil,
		},
		{
			pattern: []string{"one", "two"},
			method:  "POST",
			handler: nil,
		},
	}

	method1 := "GET"
	result1, err1 := matchMethod(routes, method1)
	if err1 != nil {
		t.Errorf("Expected no error, got '%s'", err1)
	}
	if result1.method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", result1.method)
	}

	method2 := "POST"
	result2, err2 := matchMethod(routes, method2)
	if err2 != nil {
		t.Errorf("Expected no error, got '%s'", err2)
	}
	if result2.method != "POST" {
		t.Errorf("Expected method 'POST', got '%s'", result2.method)
	}
}
