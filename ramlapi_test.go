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
	if count != 7 {
		t.Errorf("expected 7 endpoints, got %d", count)
	}

	expectedHandlers := []string{"Get", "Put", "Post", "Patch", "Delete", "Head", "NestedGet"}
	for _, h := range expectedHandlers {
		if handlerNotFound(h, endpoints) {
			t.Errorf(`expected handler name "%s", not found`, h)
		}
	}

	e1 := findEndpoint("GET", "/testapi/{foo}", endpoints)
	path := e1.Path
	if path != "/testapi/{foo}" {
		t.Errorf(`expected "/testapi/{foo}", got "%s"`, path)
	}
	if len(e1.URIParameters) != 1 {
		t.Errorf("expected 1 URIParameter, got %d", len(e1.URIParameters))
	}
	p1 := e1.URIParameters[0]
	if p1.Key != "foo" ||
		p1.Pattern != "[0-9]{5}" ||
		p1.Type != "string" ||
		p1.Required != false {
		t.Errorf("unexpected parameter values: %#v", p1)
	}

	e2 := findEndpoint("GET", "/testapi/{foo}/{bar}", endpoints)
	if len(e2.URIParameters) != 2 {
		t.Errorf("expected 2 URIParameters, got %d", len(e2.URIParameters))
	}

	if p1 != e2.URIParameters[0] {
		t.Errorf("expected endpoint to contain parameter %#v", p1)
	}

	p2 := e2.URIParameters[1]
	if p2.Key != "bar" ||
		p2.Pattern != "[a-z]{5}" ||
		p2.Type != "string" ||
		p2.Required != true {
		t.Errorf("unexpected parameter values: %#v", p2)
	}

	expectedVerbs := []string{"GET", "PUT", "POST", "PATCH", "DELETE", "HEAD"}
	for _, v := range expectedVerbs {
		if verbNotFound(v, endpoints) {
			t.Errorf(`expected verb "%s", not found`, v)
		}
	}

	// Query parameters
	queries := e1.QueryParameters
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

func queryNotFound(key string, req bool, qps []*Parameter) bool {
	for _, qp := range qps {
		if qp.Key == key && qp.Required == req {
			return false
		}
	}

	return true
}

func findEndpoint(verb, path string, eps []*Endpoint) *Endpoint {
	for _, ep := range eps {
		if ep.Verb == verb && ep.Path == path {
			return ep
		}
	}

	return nil
}
