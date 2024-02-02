package tests

import (
	"net/http"
	"pengoe/internal/router"
	"pengoe/internal/utils"
	"testing"
)

func TestAddRoute(t *testing.T) {
	// change dir to root
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	t.Errorf("Error getting current dir: %s", err)
	// }
	//
	// targetDir := filepath.Join(currentDir, "..")
	// if err := os.Chdir(targetDir); err != nil {
	// 	t.Fatal("Error changing working directory:", err)
	// }
	//
	// t.Cleanup(func() {
	// 	if err := os.Chdir(currentDir); err != nil {
	// 		t.Errorf("Error changing dir to '%s': %s", currentDir, err)
	// 	}
	// })
	//
	// t.Log(fmt.Sprintf("Current dir: %s", targetDir))

	r := router.NewRouter()

	var handler router.HandlerFunc = func(w http.ResponseWriter, r *http.Request, p map[string]string) error {
		return nil
	}

	r.AddRoute("GET", []string{"test", ":id", "helo"}, handler)

	if r.Routes[0].Method != "GET" {
		t.Errorf("Expected method 'GET', got '%s'", r.Routes[0].Method)
	}

	if utils.SliceEqual(r.Routes[0].Pattern, []string{"test", ":id", "helo"}) == false {
		t.Errorf("Expected pattern '[test :id helo]', got '%s'", r.Routes[0].Pattern)
	}

	if r.Routes[0].Handler == nil {
		t.Errorf("Expected handler address to not be nil")
	}
}

// func TestRemoveTrailingSlash(t *testing.T) {
// 	path := "/test/path/"
// 	expected := "/test/path"
//
// 	result := router.RemoveTrailingSlash(path)
//
// 	if result != expected {
// 		t.Errorf("Expected '%s', got '%s'", expected, result)
// 	}
// }
//
// func TestGetSameLengthRoutes(t *testing.T) {
// 	routes := []*router.Route{
// 		{
// 			Pattern: []string{"one", "two"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 		{
// 			Pattern: []string{"one", "two", "three"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 		{
// 			Pattern: []string{"one", "two", "other"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 		{
// 			Pattern: []string{"one", "two", "three", "four"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 	}
//
// 	path := []string{"one", "two", "three"}
//
// 	result := router.GetSameLengthRoutes(routes, path)
//
// 	if len(result) != 2 {
// 		t.Errorf("Expected 2 routes, got %d", len(result))
// 	}
//
// 	for _, route := range result {
// 		if len(route.Pattern) != len(path) {
// 			t.Errorf("Expected route pattern length to be %d, got %d", len(path), len(route.Pattern))
// 		}
// 	}
// }
//
// func TestMatchRoutes(t *testing.T) {
// 	routes := []*router.Route{
// 		{
// 			Pattern: []string{"one", "two"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 		{
// 			Pattern: []string{"one", ":var"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 	}
//
// 	path1 := []string{"one", "two"}
// 	result1, err1 := router.MatchRoutes(routes, path1)
// 	if err1 != nil {
// 		t.Errorf("Expected no error, got '%s'", err1)
// 	}
// 	if len(result1) != 1 {
// 		t.Errorf("Expected 1 route, got %d", len(result1))
// 	}
// 	if !utils.SliceEqual(result1[0].Pattern, []string{"one", "two"}) {
// 		t.Errorf("Expected route pattern to be '[one two]', got '%s'", result1[0].Pattern)
// 	}
//
// 	path2 := []string{"one", "three"}
// 	result2, err2 := router.MatchRoutes(routes, path2)
// 	if err2 != nil {
// 		t.Errorf("Expected no error, got '%s'", err2)
// 	}
// 	if len(result2) != 1 {
// 		t.Errorf("Expected 1 route, got %d", len(result2))
// 	}
// 	if !utils.SliceEqual(result2[0].Pattern, []string{"one", ":var"}) {
// 		t.Errorf("Expected route pattern to be '[one :var]', got '%s'", result2[0].Pattern)
// 	}
//
// 	path3 := []string{"two", "three"}
// 	result3, err3 := router.MatchRoutes(routes, path3)
// 	if err3 == nil {
// 		t.Errorf("Expected error, got nil")
// 	}
// 	if len(result3) != 0 {
// 		t.Errorf("Expected 0 routes, got %d", len(result3))
// 	}
// }
//
// func TestMatchMethod(t *testing.T) {
// 	routes := []*router.Route{
// 		{
// 			Pattern: []string{"one", "two"},
// 			Method:  "GET",
// 			Handler: nil,
// 		},
// 		{
// 			Pattern: []string{"one", "two"},
// 			Method:  "POST",
// 			Handler: nil,
// 		},
// 	}
//
// 	method1 := "GET"
// 	result1, err1 := router.MatchMethod(routes, method1)
// 	if err1 != nil {
// 		t.Errorf("Expected no error, got '%s'", err1)
// 	}
// 	if result1.Method != "GET" {
// 		t.Errorf("Expected method 'GET', got '%s'", result1.Method)
// 	}
//
// 	method2 := "POST"
// 	result2, err2 := router.MatchMethod(routes, method2)
// 	if err2 != nil {
// 		t.Errorf("Expected no error, got '%s'", err2)
// 	}
// 	if result2.Method != "POST" {
// 		t.Errorf("Expected method 'POST', got '%s'", result2.Method)
// 	}
// }
