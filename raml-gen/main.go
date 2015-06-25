package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/EconomistDigitalSolutions/ramlapi"
	"github.com/buddhamagnet/raml"
)

var (
	ramlFile string
	genFile  string
)

func init() {
	flag.StringVar(&ramlFile, "ramlfile", "api.raml", "RAML file to parse")
	flag.StringVar(&genFile, "genfile", "handlers_gen.go", "Filename to use for output")
}

func main() {
	flag.Parse()
	api, err := ramlapi.ProcessRAML(ramlFile)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Processing API spec for", ramlFile)
	generate(api, genFile)
	log.Println("Created handlers in ", genFile)
}

// Generate handler functions based on an API definition.
func generate(api *raml.APIDefinition, genFile string) {

	f, err := os.Create(genFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Write the header - import statements and root handler.
	f.WriteString(handlerHead)
	// Start the route map (string to handler).
	f.WriteString(mapStart)
	// Add the route map entries.
	e := template.Must(template.New("mapEntry").Parse(mapEntry))
	for name, resource := range api.Resources {
		generateMap("", name, &resource, e, f)
	}
	// Close the route map.
	f.WriteString(mapEnd)
	// Now add the HTTP handlers.
	t := template.Must(template.New("handlerText").Parse(handlerText))
	for name, resource := range api.Resources {
		generateResource("", name, &resource, t, f)
	}
	format(f)
}

// format runs go fmt on a file.
func format(f *os.File) {
	// Run go fmt on the file.
	cmd := exec.Command("go", "fmt")
	cmd.Stdin = f
	_ = cmd.Run()
}

// generateResource creates a handler struct from an API resource
// and executes the associated template.
func generateResource(parent, name string, resource *raml.Resource, t *template.Template, f *os.File) string {
	var resourcepath = parent + name
	type HandlerInfo struct {
		Name, Verb, Path, Doc string
	}

	for verb, n := range ramlapi.ResourceVerbs(resource) {
		if n == "" {
			log.Fatalf("no handler name specified for %s via %s\n", resourcepath, verb)
		}
		err := t.Execute(f, HandlerInfo{n, verb, resourcepath, resource.Description})
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		return generateResource(resourcepath, nestname, nested, t, f)
	}
	return resourcepath
}

// generateMap builds a map of string labels to handler funcs - this is
// used by the calling code to link the display name strings that come
// from the RAML file to handler funcs in the client code.
func generateMap(parent, name string, resource *raml.Resource, e *template.Template, f *os.File) {
	var resourcepath = parent + name
	type RouteMapEntry struct {
		Name, Struct string
	}

	for verb, n := range ramlapi.ResourceVerbs(resource) {
		if n == "" {
			log.Fatalf("no handler name specified for %s via %s\n", resourcepath, verb)
		}
		err := e.Execute(f, RouteMapEntry{n, n})
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		generateMap(resourcepath, nestname, nested, e, f)
	}
}