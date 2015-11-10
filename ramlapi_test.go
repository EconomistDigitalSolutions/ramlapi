package ramlapi_test

import (
	"testing"

	. "github.com/EconomistDigitalSolutions/ramlapi"
)

var endpoints []*Endpoint

func testFunc(ep *Endpoint) {
	endpoints = append(endpoints, ep)
}

func TestAPI(t *testing.T) {
	api, _ := Process("fixtures/valid.raml")
	Build(api, testFunc)

	count := len(endpoints)
	if count != 6 {
		t.Errorf("expected 6 endpoints, got %d", count)
	}

	expectedHandlers := []string{"Get", "Put", "Post", "Patch", "Delete", "Head"}
	for _, h := range expectedHandlers {
		if handlerNotFound(h, endpoints) {
			t.Errorf(`expected handler name "%s", not found`, h)
		}
	}

	path := endpoints[0].Path
	if path != "/testapi" {
		t.Errorf(`expected "/testapi", got "%s"`, path)
	}

	expectedVerbs := []string{"GET", "PUT", "POST", "PATCH", "DELETE", "HEAD"}
	for _, v := range expectedVerbs {
		if verbNotFound(v, endpoints) {
			t.Errorf(`expected verb "%s", not found`, v)
		}
	}

	get := findGetEndpoint(endpoints)
	queries := get.QueryParameters
	if len(queries) != 2 {
		t.Errorf("expected 1 query string parameter, got %d", len(queries))
	}

	// query map {key: required}
	expectedQueries := map[string]bool{"Country": true, "City": false}
	for key, req := range expectedQueries {
		if queryNotFound(key, req, queries) {
			t.Errorf("expected parameter: %s, required: %t, none found", key, req)
		}
	}
}

func handlerNotFound(handler string, eps []*Endpoint) bool {
	for _, ep := range eps {
		if ep.Handler == handler {
			return false
		}
	}

	return true
}

func verbNotFound(verb string, eps []*Endpoint) bool {
	for _, ep := range eps {
		if ep.Verb == verb {
			return false
		}
	}

	return true
}

func queryNotFound(key string, req bool, qps []*QueryParameter) bool {
	for _, qp := range qps {
		if qp.Key == key && qp.Required == req {
			return false
		}
	}

	return true
}

func findGetEndpoint(eps []*Endpoint) *Endpoint {
	for _, ep := range eps {
		if ep.Verb == "GET" {
			return ep
		}
	}

	return nil
}
