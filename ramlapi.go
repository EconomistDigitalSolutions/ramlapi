package ramlapi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/buddhamagnet/raml"
)

var reg = regexp.MustCompile("[^A-Za-z0-9]+")

// QueryParameter represents a URL query parameter.
type QueryParameter struct {
	Key      string
	Type     string
	Pattern  string
	Required bool
}

// Endpoint describes an API endpoint.
type Endpoint struct {
	Verb            string
	Handler         string
	Path            string
	QueryParameters []*QueryParameter
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("verb: %s handler: %s path:%s\n", e.Verb, e.Handler, e.Path)
}

func (e *Endpoint) setQueryParameters(method *raml.Method) {
	for _, res := range method.QueryParameters {
		q := &QueryParameter{
			Key:      res.Name,
			Type:     res.Type,
			Required: res.Required,
		}
		if res.Pattern != nil {
			q.Pattern = *res.Pattern
		}
		e.QueryParameters = append(e.QueryParameters, q)
	}
}

// EndpointSet is a set of API endpoints.
type EndpointSet struct {
	Endpoints []*Endpoint
}

func (s *EndpointSet) addEndpoint(verb string, method *raml.Method) {
	if method != nil {
		ep := &Endpoint{
			Verb:    verb,
			Handler: reg.ReplaceAllString(strings.Title(method.DisplayName), ""),
		}
		ep.setQueryParameters(method)

		s.Endpoints = append(s.Endpoints, ep)
	}
}

// ProcessRAML processes a RAML file and returns an API definition.
func ProcessRAML(ramlFile string) (*raml.APIDefinition, error) {
	routes, err := raml.ParseFile(ramlFile)
	if err != nil {
		return nil, fmt.Errorf("Failed parsing RAML file: %s\n", err.Error())
	}
	return routes, nil
}

// processResource recursively process resources and their nested children
// and returns the path so far for the children. The function takes a routerFunc
// as an argument that is invoked with the verb, resource path and handler as
// the resources are processed, so the calling code can use pat, mux, httprouter
// or whatever router they desire and we don't need to know about it.
func processResource(parent, name string, resource *raml.Resource, routerFunc func(s *Endpoint)) string {
	var path = parent + name

	s := new(EndpointSet)
	s.addEndpoint("GET", resource.Get)
	s.addEndpoint("POST", resource.Post)
	s.addEndpoint("PUT", resource.Put)
	s.addEndpoint("PATCH", resource.Patch)
	s.addEndpoint("HEAD", resource.Head)
	s.addEndpoint("DELETE", resource.Delete)

	for _, ep := range s.Endpoints {
		ep.Path = path
		routerFunc(ep)
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		return processResource(path, nestname, nested, routerFunc)
	}

	return path
}

// Build takes a RAML API definition, a router and a routing map,
// and wires them all together.
func Build(api *raml.APIDefinition, routerFunc func(s *Endpoint)) {
	for name, resource := range api.Resources {
		processResource("", name, &resource, routerFunc)
	}
}
