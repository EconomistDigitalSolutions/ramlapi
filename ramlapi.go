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

	for verb, handler := range ResourceVerbs(resource) {
		log.Println("--- " + verb)
		data := map[string]string{
			"verb":    verb,
			"path":    resourcepath,
			"handler": handler,
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
func ResourceVerbs(resource *raml.Resource) map[string]string {
	var verbs = make(map[string]string)

	if resource.Get != nil {
		verbs["GET"] = resource.Get.DisplayName
	}
	if resource.Post != nil {
		verbs["POST"] = resource.Post.DisplayName
	}
	if resource.Put != nil {
		verbs["PUT"] = resource.Put.DisplayName
	}
	if resource.Patch != nil {
		verbs["PATCH"] = resource.Patch.DisplayName
	}
	if resource.Head != nil {
		verbs["HEAD"] = resource.Head.DisplayName
	}
	if resource.Delete != nil {
		verbs["DELETE"] = resource.Delete.DisplayName
	}

	return verbs
}
