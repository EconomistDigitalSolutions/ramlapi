package ramlapi_test

import (
	"fmt"
	"testing"

	. "github.com/EconomistDigitalSolutions/ramlapi"
	"github.com/buddhamagnet/raml"
)

var (
	endpoints []*Endpoint
	patt1     string
	patt2     string
)

func init() {
	patt1 = "[a-z]"
	patt2 = "[0-9]"
}

func testFunc(ep *Endpoint) {
	endpoints = append(endpoints, ep)
}

// TestData used to test raml.APIDefinition to ramlapi.Endpoints. The order
// of expected methods is important (GET, POST, PUT, PATCH, HEAD, DELETE).
var TestData = []struct {
	api      *raml.APIDefinition
	expected []map[string]interface{}
}{
	{
		// test simple endpoint
		&raml.APIDefinition{
			Resources: map[string]raml.Resource{
				"/test": raml.Resource{
					Post: &raml.Method{
						Name:        "POST",
						DisplayName: "Post me",
					},
					Get: &raml.Method{
						Name:        "GET",
						DisplayName: "Get me",
					},
				},
			},
		},
		[]map[string]interface{}{
			{"verb": "GET", "handler": "GetMe", "path": "/test"},
			{"verb": "POST", "handler": "PostMe", "path": "/test"},
		},
	},
	{
		// test URI parameters with nested resource
		&raml.APIDefinition{
			Resources: map[string]raml.Resource{
				"/{foo}": raml.Resource{
					Get: &raml.Method{
						Name:        "GET",
						DisplayName: "Get me",
					},
					UriParameters: map[string]raml.NamedParameter{
						"foo": raml.NamedParameter{
							Pattern: &patt1,
						},
					},
					Nested: map[string]*raml.Resource{
						"/{bar}": &raml.Resource{
							Get: &raml.Method{
								Name:        "GET",
								DisplayName: "Nested get",
							},
							UriParameters: map[string]raml.NamedParameter{
								"bar": raml.NamedParameter{
									Pattern: &patt2,
								},
							},
						},
					},
				},
			},
		},
		[]map[string]interface{}{
			{
				"verb":    "GET",
				"handler": "GetMe",
				"path":    "/{foo}",
				"uri_params": []map[string]string{
					{
						"key":     "foo",
						"pattern": "[a-z]",
					},
				},
			},
			{
				"verb":    "GET",
				"handler": "NestedGet",
				"path":    "/{foo}/{bar}",
				"uri_params": []map[string]string{
					{
						"key":     "foo",
						"pattern": "[a-z]",
					},
					{
						"key":     "bar",
						"pattern": "[0-9]",
					},
				},
			},
		},
	},
	{
		// test query parameters
		&raml.APIDefinition{
			Resources: map[string]raml.Resource{
				"/query": raml.Resource{
					Get: &raml.Method{
						Name:        "GET",
						DisplayName: "Get me",
						QueryParameters: map[string]raml.NamedParameter{
							"foo": raml.NamedParameter{
								Pattern:  &patt1,
								Required: true,
							},
							"bar": raml.NamedParameter{
								Pattern:  &patt2,
								Required: false,
							},
						},
					},
				},
			},
		},
		[]map[string]interface{}{
			{
				"verb":    "GET",
				"handler": "GetMe",
				"path":    "/query",
				"query_params": []map[string]string{
					{
						"key":      "foo",
						"pattern":  "[a-z]",
						"required": "true",
					},
					{
						"key":      "bar",
						"pattern":  "[0-9]",
						"required": "false",
					},
				},
			},
		},
	},
}

func TestProcess(t *testing.T) {
	_, err := Process("fixtures/valid.raml")
	if err != nil {
		t.Error("could not process valid RAML file")
	}
}

func TestEndpoints(t *testing.T) {
	for _, data := range TestData {
		Build(data.api, testFunc)
		if !checkEndpoints(t, data.expected, endpoints) {
			t.Errorf("expected endpoint with: %s", data.expected)
		}
		endpoints = make([]*Endpoint, 0)
	}
}

func checkEndpoints(t *testing.T, exp []map[string]interface{}, got []*Endpoint) bool {
	var foundHandler, foundPath, foundVerb bool
	var found int

	if len(exp) != len(got) {
		t.Errorf("expected %d endpoints, got %d", len(exp), len(got))
	}

	for _, e := range exp {
		for _, ep := range got {

			foundHandler = true
			if ep.Handler != e["handler"] {
				foundHandler = false
			}

			foundPath = true
			if ep.Path != e["path"] {
				foundPath = false
			}

			foundVerb = true
			if ep.Verb != e["verb"] {
				foundVerb = false
			}

			if foundHandler && foundPath && foundVerb {
				found += 1
				if uParams, ok := e["uri_params"]; ok {
					u := uParams.([]map[string]string)
					if !checkParameters(t, u, ep.URIParameters) {
						t.Errorf("expected uri parameters: %s", u)
					}
				}
				if qParams, ok := e["query_params"]; ok {
					q := qParams.([]map[string]string)
					if !checkParameters(t, q, ep.QueryParameters) {
						t.Errorf("expected query parameters: %s", q)
					}
				}
			}
		}
	}

	return found == len(exp)
}

func checkParameters(t *testing.T, exp []map[string]string, got []*Parameter) bool {
	var foundKey, foundPatt, foundReq bool
	var found int

	if len(got) != len(exp) {
		t.Errorf("expected %d parameters, got %d", len(exp), len(got))
	}

	for _, expParam := range exp {
		for _, gotParam := range got {
			foundKey = true
			if key, ok := expParam["key"]; ok {
				if key != gotParam.Key {
					foundKey = false
				}
			}

			foundPatt = true
			if patt, ok := expParam["pattern"]; ok {
				if patt != gotParam.Pattern {
					foundPatt = false
				}
			}

			foundReq = true
			if req, ok := expParam["required"]; ok {
				// convert bool to string
				gotReq := fmt.Sprintf("%t", gotParam.Required)
				if req != gotReq {
					foundReq = false
				}
			}

			if foundKey && foundPatt && foundReq {
				found += 1
			}
		}
	}

	return found == len(exp)
}
