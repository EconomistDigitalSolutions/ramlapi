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

// Parameter is a path or query string parameter.
type Parameter struct {
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
	URIParameters   []*Parameter
	QueryParameters []*Parameter
}

// String returns the string representation of an Endpoint.
func (e *Endpoint) String() string {
	return fmt.Sprintf("verb: %s handler: %s path:%s\n", e.Verb, e.Handler, e.Path)
}

func (e *Endpoint) setQueryParameters(method *raml.Method) {
	for name, param := range method.QueryParameters {
		e.QueryParameters = append(e.QueryParameters, newParam(name, &param))
	}
}

// Build takes a RAML API definition, a router and a routing map,
// and wires them all together.
func Build(api *raml.APIDefinition, routerFunc func(s *Endpoint)) error {
	for name, resource := range api.Resources {
		var resourceParams []*Parameter
		err := processResource("", name, &resource, resourceParams, routerFunc)
		if err != nil {
			return err
		}
	}

	return nil
}

// Process processes a RAML file and returns an API definition.
func Process(file string) (*raml.APIDefinition, error) {
	routes, err := raml.ParseFile(file)
	if err != nil {
		return nil, fmt.Errorf("Failed parsing RAML file: %s\n", err.Error())
	}
	return routes, nil
}

// Variableize normalises RAML display names.
func Variableize(s string) string {
	return vizer.ReplaceAllString(strings.Title(s), "")
}

func newParam(name string, param *raml.NamedParameter) *Parameter {
	p := &Parameter{
		Key:      name,
		Type:     param.Type,
		Required: param.Required,
	}
	if param.Pattern != nil {
		p.Pattern = *param.Pattern
	}

	return p
}

func appendEndpoint(s []*Endpoint, method *raml.Method, params []*Parameter) ([]*Endpoint, error) {
	if method.DisplayName == "" {
		return s, errors.New("DisplayName property not set in RAML method")
	}

	if method != nil {
		ep := &Endpoint{
			Verb:        method.Name,
			Handler:     Variableize(method.DisplayName),
			Description: method.Description,
		}
		// set query parameters
		ep.setQueryParameters(method)
		// set uri parameters
		for _, param := range params {
			ep.URIParameters = append(ep.URIParameters, param)
		}
		s = append(s, ep)
	}

	return s, nil
}

// processResource recursively process resources and their nested children
// and returns the path so far for the children. The function takes a routerFunc
// as an argument that is invoked with the verb, resource path and handler as
// the resources are processed, so the calling code can use pat, mux, httprouter
// or whatever router they desire and we don't need to know about it.
func processResource(parent, name string, resource *raml.Resource, params []*Parameter, routerFunc func(s *Endpoint)) error {
	var path = parent + name
	var err error
	for name, param := range resource.UriParameters {
		params = append(params, newParam(name, &param))
	}

	s := make([]*Endpoint, 0, 6)
	for _, m := range resource.Methods() {
		s, err = appendEndpoint(s, m, params)
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
		return processResource(path, nestname, nested, params, routerFunc)
	}

	return nil
}
