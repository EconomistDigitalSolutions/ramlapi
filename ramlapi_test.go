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
	api, _ := ProcessRAML("fixtures/valid.raml")
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

	if !queries[1].Required {
		t.Error("expected query to be required")
	}

	if queries[0].Key != "country" {
		// TODO: add NamedParameter.Name to raml package so this test doesn't fail
		t.Errorf(`expected paramater key to be "country", got "%s"`, queries[0].Key)
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

func findGetEndpoint(eps []*Endpoint) *Endpoint {
	for _, ep := range eps {
		if ep.Verb == "GET" {
			return ep
		}
	}

	return nil
}
