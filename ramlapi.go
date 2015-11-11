package ramlapi

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/buddhamagnet/raml"
)

var vizer *regexp.Regexp

func init() {
	vizer = regexp.MustCompile("[^A-Za-z0-9]+")
}

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
	Description     string
	QueryParameters []*QueryParameter
}

// String returns the string representation of an Endpoint.
func (e *Endpoint) String() string {
	return fmt.Sprintf("verb: %s handler: %s path:%s\n", e.Verb, e.Handler, e.Path)
}

func (e *Endpoint) setQueryParameters(method *raml.Method) {
	for name, res := range method.QueryParameters {
		q := &QueryParameter{
			Key:      name,
			Type:     res.Type,
			Required: res.Required,
		}
		if res.Pattern != nil {
			q.Pattern = *res.Pattern
		}
		e.QueryParameters = append(e.QueryParameters, q)
	}
}

// Process processes a RAML file and returns an API definition.
func Process(file string) (*raml.APIDefinition, error) {
	routes, err := raml.ParseFile(file)
	if err != nil {
		return nil, fmt.Errorf("Failed parsing RAML file: %s\n", err.Error())
	}
	return routes, nil
}

// Build takes a RAML API definition, a router and a routing map,
// and wires them all together.
func Build(api *raml.APIDefinition, routerFunc func(s *Endpoint)) error {
	for name, resource := range api.Resources {
		err := processResource("", name, &resource, routerFunc)
		if err != nil {
			return err
		}
	}

	return nil
}

func appendEndpoint(s []*Endpoint, method *raml.Method) ([]*Endpoint, error) {
	if method.DisplayName == "" {
		return s, errors.New("DisplayName property not set in RAML method")
	}

	if method != nil {
		ep := &Endpoint{
			Verb:        method.Name,
			Handler:     Variableize(method.DisplayName),
			Description: method.Description,
		}
		ep.setQueryParameters(method)

		s = append(s, ep)
	}

	return s, nil
}

// processResource recursively process resources and their nested children
// and returns the path so far for the children. The function takes a routerFunc
// as an argument that is invoked with the verb, resource path and handler as
// the resources are processed, so the calling code can use pat, mux, httprouter
// or whatever router they desire and we don't need to know about it.
func processResource(parent, name string, resource *raml.Resource, routerFunc func(s *Endpoint)) error {
	var path = parent + name
	var err error

	s := make([]*Endpoint, 0, 6)
	for _, details := range resource.Methods() {
		s, err = appendEndpoint(s, details)
		if err != nil {
			return err
		}
	}

	for _, ep := range s {
		ep.Path = path
		routerFunc(ep)
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		return processResource(path, nestname, nested, routerFunc)
	}

	return nil
}

// Variableize normalises RAML display names.
func Variableize(s string) string {
	return vizer.ReplaceAllString(strings.Title(s), "")
}
