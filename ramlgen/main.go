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

// RouteMapEntry represents an entry in a route map.
type RouteMapEntry struct {
	Name, Struct string
}

// HandlerInfo contains handler information.
type HandlerInfo struct {
	Name, Verb, Path, Doc string
}

func init() {
	flag.StringVar(&ramlFile, "ramlfile", "api.raml", "RAML file to parse")
	flag.StringVar(&genFile, "genfile", "handlers_gen.go", "Filename to use for output")
}

func main() {
	flag.Parse()
	api, err := ramlapi.Process(ramlFile)
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
	path := parent + name

	for _, method := range resource.Methods() {
		err := t.Execute(f, HandlerInfo{ramlapi.Variableize(method.DisplayName), method.Name, path, method.Description})
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		return generateResource(path, nestname, nested, t, f)
	}
	return path
}

// generateMap builds a map of string labels to handler funcs - this is
// used by the calling code to link the display name strings that come
// from the RAML file to handler funcs in the client code.
func generateMap(parent, name string, resource *raml.Resource, e *template.Template, f *os.File) {
	path := parent + name

	for _, method := range resource.Methods() {
		name := ramlapi.Variableize(method.DisplayName)
		err := e.Execute(f, RouteMapEntry{name, name})
		if err != nil {
			log.Println("executing template:", err)
		}
	}

	// Get all children.
	for nestname, nested := range resource.Nested {
		generateMap(path, nestname, nested, e, f)
	}
}
