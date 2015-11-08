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

	if endpoints[0].Handler != "Get" {
		t.Errorf(`expected handler name "Get", got "%s"`, endpoints[0].Handler)
	}

	path := endpoints[0].Path
	if path != "/testapi" {
		t.Errorf(`expected "/testapi", got "%s"`, path)
	}

	verb := endpoints[0].Verb
	if verb != "GET" {
		// TODO: add Method.Name to raml package so this test doesn't fail
		t.Errorf(`expected "GET", got "%s"`, verb)
	}

	queries := endpoints[0].QueryParameters
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
