package ramlapi

import (
	"testing"
)

var router RouterMock

type RouterMock struct {
	Set *EndpointSet
}

func (r *RouterMock) Consume(s *EndpointSet) {
	r.Set = s
}

func routerFunc(s *EndpointSet) {
	router.Consume(s)
}

func TestAPI(t *testing.T) {
	// TODO: consider fixtrue data as a bunch of structs
	api, _ := ProcessRAML("fixtures/valid.raml")
	Build(api, routerFunc)

	count := len(router.Set.Endpoints)
	if count != 6 {
		t.Errorf("expected 6 endpoints, got %d", count)
	}

	if router.Set.Endpoints[0].Handler != "GetTestEndpoint" {
		t.Errorf(`expected handler name "GetTestEndpoint", got "%s"`, router.Set.Endpoints[0].Handler)
	}

	path := router.Set.Endpoints[0].Path
	if path != "/testapi" {
		t.Errorf(`expected "/testapi", got "%s"`, path)
	}

	verb := router.Set.Endpoints[0].Verb
	if verb != "GET" {
		// TODO: add Method.Name to raml package so this test doesn't fail
		t.Errorf(`expected "GET", got "%s"`, verb)
	}

	queries := router.Set.Endpoints[0].QueryParameters
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
