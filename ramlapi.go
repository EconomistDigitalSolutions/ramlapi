package ramlapi

import (
	"fmt"
	"log"
	"os"

	"github.com/buddhamagnet/raml"
)

// Build takes a RAML API definition, a router and a routing map,
// and wires them all together.
func Build(api *raml.APIDefinition, fun func(data map[string]string)) {
	for name, res := range api.Resources {
		processResource("", name, &res, fun)
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

// ResourceVerbs assembles resource method types into a
// map of verbs to handler names.
func ResourceVerbs(res *raml.Resource) map[string]map[string]string {
	var verbs = make(map[string]map[string]string)

	if res.Get != nil {
		verbs["GET"] = map[string]string{
			"handler": res.Get.DisplayName,
		}
		if len(res.Get.QueryParameters) >= 1 {
			for _, value := range res.Get.QueryParameters {
				verbs["GET"]["query"] = value.DisplayName
				verbs["GET"]["query_type"] = value.Type
				verbs["GET"]["query_description"] = value.Description
				verbs["GET"]["query_example"] = value.Example
				verbs["GET"]["query_pattern"] = *value.Pattern
			}
		}
	}

	mappings := map[string]*raml.Method{
		"POST":   res.Post,
		"PUT":    res.Put,
		"PATCH":  res.Patch,
		"DELETE": res.Delete,
		"HEAD":   res.Head,
	}

	for verb, meth := range mappings {
		if meth != nil {
			verbs[verb] = map[string]string{
				"handler": meth.DisplayName,
			}
		}

	}

	return verbs
}

// processResource recursively process resources and their nested children
// and returns the path so far for the children. The function takes a routerFunc
// as an argument that is invoked with the verb, resource path and handler as
// the resources are processed, so the calling code can use pat, mux, httprouter
// or whatever router they desire and we don't need to know about it.
func processResource(parent, name string, res *raml.Resource, fun func(data map[string]string)) string {
	path := parent + name

	logg("processing", name, "resource")
	logg("path: ", path)

	for verb, details := range ResourceVerbs(res) {
		logg("--- " + verb)
		data := map[string]string{
			"verb":    verb,
			"path":    path,
			"handler": details["handler"],
		}
		if details["query"] != "" {
			data["query"] = details["query"]
			data["query_type"] = details["query_type"]
			data["query_description"] = details["query_description"]
			data["query_example"] = details["query_example"]
			data["query_pattern"] = details["query_pattern"]
		}
		fun(data)
	}

	// Get all children.
	for name, n := range res.Nested {
		return processResource(path, name, n, fun)
	}
	return path
}

func logg(args ...string) {
	if os.Getenv("RAMLAPI_LOGGING") == "" {
		log.Println(args)
	}
}
