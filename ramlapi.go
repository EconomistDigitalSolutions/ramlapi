package ramlapi

import (
	"fmt"
	"log"

	"github.com/buddhamagnet/raml"
)

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
func processResource(parent, name string, resource *raml.Resource, routerFunc func(data map[string]string)) string {

	var resourcepath = parent + name
	log.Println("processing", name, "resource")
	log.Println("path: ", resourcepath)

	for verb, details := range ResourceVerbs(resource) {
		log.Println("--- " + verb)
		data := map[string]string{
			"verb":              verb,
			"path":              resourcepath,
			"handler":           details["handler"],
			"query":             details["query"],
			"query_type":        details["query_type"],
			"query_description": details["query_description"],
			"query_example":     details["query_example"],
			"query_pattern":     details["query_pattern"],
		}
		routerFunc(data)
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		return processResource(resourcepath, nestname, nested, routerFunc)
	}
	return resourcepath
}

// Build takes a RAML API definition, a router and a routing map,
// and wires them all together.
func Build(api *raml.APIDefinition, routerFunc func(data map[string]string)) {
	for name, resource := range api.Resources {
		processResource("", name, &resource, routerFunc)
	}
}

// ResourceVerbs assembles resource method types into a
// map of verbs to handler names.
func ResourceVerbs(resource *raml.Resource) map[string]map[string]string {
	var verbs = make(map[string]map[string]string)

	if resource.Get != nil {
		verbs["GET"] = map[string]string{
			"handler": resource.Get.DisplayName,
		}
		if len(resource.Get.QueryParameters) >= 1 {
			for _, value := range resource.Get.QueryParameters {
				verbs["GET"]["query"] = value.DisplayName
				verbs["GET"]["query_type"] = value.Type
				verbs["GET"]["query_description"] = value.Description
				verbs["GET"]["query_example"] = value.Example
				verbs["GET"]["query_pattern"] = *value.Pattern
			}
		}
	}
	if resource.Post != nil {
		verbs["POST"] = map[string]string{
			"handler": resource.Post.DisplayName,
		}
	}
	if resource.Put != nil {
		verbs["PUT"] = map[string]string{
			"handler": resource.Put.DisplayName,
		}
	}
	if resource.Patch != nil {
		verbs["PATCH"] = map[string]string{
			"handler": resource.Patch.DisplayName,
		}
	}
	if resource.Head != nil {
		verbs["HEAD"] = map[string]string{
			"handler": resource.Head.DisplayName,
		}
	}
	if resource.Delete != nil {
		verbs["DELETE"] = map[string]string{
			"handler": resource.Delete.DisplayName,
		}
	}

	return verbs
}
