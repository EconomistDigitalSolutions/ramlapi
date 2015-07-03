package ramlapi

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var handlers = map[string]http.HandlerFunc{
	"GetMe":    GetMe,
	"PostMe":   PostMe,
	"PutMe":    PutMe,
	"PatchMe":  PatchMe,
	"HeadMe":   HeadMe,
	"DeleteMe": DeleteMe,
}

var router = mux.NewRouter().StrictSlash(true)
var routerMock RouterMock

type RouterMock struct {
	routes map[string]map[string]string
}

func (router *RouterMock) Consume(data map[string]string) {
	router.routes = make(map[string]map[string]string)
	router.routes[data["path"]] = map[string]string{
		data["verb"]: data["handler"],
	}
}

func (router *RouterMock) String() {
	for route, data := range router.routes {
		fmt.Println(route)
		for k, v := range data {
			fmt.Println(k, v)
		}
	}
}

func (router *RouterMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.String()
	if _, found := router.routes[path]; found {
		if _, found := router.routes[path][r.Method]; found {
			handler := handlers[(router.routes[path][r.Method])]
			handler(w, r)
		}
	}
}

func routerFunc(data map[string]string) {
	router.
		Methods(data["verb"]).
		Path(data["path"]).
		Handler(handlers[data["handler"]])
}

func routerFuncMock(data map[string]string) {
	routerMock.Consume(data)
}

func buildAPI() {
	api, _ := ProcessRAML("fixtures/valid.raml")
	Build(api, routerFunc)
}

func buildAPIMock() {
	api, _ := ProcessRAML("fixtures/valid.raml")
	Build(api, routerFuncMock)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("GETME")
	w.Write([]byte("GetMe"))
}
func PostMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("POSTME")
	w.Write([]byte("PostMe"))
}

func PutMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("PUTME")
	w.Write([]byte("PutMe"))
}

func PatchMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("PATCHME")
	w.Write([]byte("PatchMe"))
}

func HeadMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("HEADME")
	w.Write([]byte("HeadMe"))
}

func DeleteMe(w http.ResponseWriter, r *http.Request) {
	log.Fatal("DELETEME")
	w.Write([]byte("DeleteMe"))
}

func TestMissingRaml(t *testing.T) {
	_, err := ProcessRAML("fixtures/missing.raml")
	if err == nil {
		t.Fatal("Expected error with missing RAML file")
	}
}

func TestInvalidRaml(t *testing.T) {
	_, err := ProcessRAML("fixtures/invalid.raml")
	if err == nil {
		t.Fatal("Expected error with invalid RAML file")
	}
}

func TestValidRaml(t *testing.T) {
	_, err := ProcessRAML("fixtures/valid.raml")
	if err != nil {
		t.Fatalf("Expected good response with valid RAML file, got %v\n", err)
	}
}

func TestValidRamlGetAssignments(t *testing.T) {
	// Build the API and assign handlers.
	buildAPI()
	// Cycle through the map and dispatch the appropriate
	// HTTP requests to each one.
	for name := range handlers {

		matcher := regexp.MustCompile("^(Get|Post|Put|Patch|Head|Delete)")
		match := matcher.FindSubmatch([]byte(name))
		req, err := http.NewRequest(strings.ToUpper(string(match[0])), "/testapi", nil)

		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		// We need to send this to the mux to ensure we are testing the
		// router interface i.e. the handlers have been assigned when the
		// API was built.
		router.ServeHTTP(res, req)

		// Now make sure every handler returns with a 200 OK and the
		// correct response body.
		if res.Code != 200 {
			t.Fatalf("Expected a 200 response from %s, got %d\n", name, res.Code)
		}
		if res.Body.String() != name {
			t.Fatalf("Expected to get %s response from %s, got %s\n", name, name, res.Body.String())
		}
	}
}

func TestValidRamlGetAssignmentsMock(t *testing.T) {
	routerMock = RouterMock{}
	// Build the API and assign handlers.
	buildAPIMock()
	// Cycle through the map and dispatch the appropriate
	// HTTP requests to each one.
	for name := range handlers {

		matcher := regexp.MustCompile("^(Get|Post|Put|Patch|Head|Delete)")
		match := matcher.FindSubmatch([]byte(name))
		req, err := http.NewRequest(strings.ToUpper(string(match[0])), "/testapi", nil)

		if err != nil {
			t.Fatal(err)
		}

		res := httptest.NewRecorder()
		// We need to send this to the mux to ensure we are testing the
		// router interface i.e. the handlers have been assigned when the
		// API was built.
		routerMock.ServeHTTP(res, req)

		// Now make sure every handler returns with a 200 OK and the
		// correct response body.
		if res.Code != 200 {
			t.Fatalf("Expected a 200 response from %s, got %d\n", name, res.Code)
		}
		if res.Body.String() != name {
			t.Fatalf("Expected to get %s response from %s, got %s\n", name, name, res.Body.String())
		}
	}
}
