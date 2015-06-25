package ramlapi

import (
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

func routerFunc(data map[string]string) {
	router.
		Methods(data["verb"]).
		Path(data["path"]).
		Handler(handlers[data["handler"]])
}

func buildAPI() {
	api, _ := ProcessRAML("fixtures/valid.raml")
	Build(api, routerFunc)
}

func GetMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("GetMe"))
}
func PostMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PostMe"))
}

func PutMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PutMe"))
}

func PatchMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PatchMe"))
}

func HeadMe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HeadMe"))
}

func DeleteMe(w http.ResponseWriter, r *http.Request) {
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
