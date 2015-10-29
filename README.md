This package is designed to work with [RAML](http://raml.org) files and build the appropriate handlers and routing rules in Go applications.

[![GoDoc](https://godoc.org/github.com/EconomistDigitalSolutions/ramlapi?status.svg)](https://godoc.org/github.com/EconomistDigitalSolutions/ramlapi)

The ramlapi codebase contains two packages:

* Ramlapi - used to parse a RAML file and wire it up to a router.
* Ramlgen - used to parse a RAML file and write a set of HTTP handlers.

#### HOW TO RAML-GEN

1. Build your API design in RAML.
2. Inside your Go project, run `raml-gen --ramlfile=<file>`.
3. Copy the resulting `handlers_gen.go` file to the correct location.

You now have a set of HTTP handlers built from your RAML specification.

#### HOW TO RAMLAPI

The ramlapi package makes no assumptions about your choice of router as the
method to wire up the API takes a function provided by your code and
passes details of the API back to that function on each resource defined
in the RAML file. The router can then hook the data up however it likes.

#### EXAMPLES

##### GORILLA MUX

```go
router = mux.NewRouter().StrictSlash(true)
api, err := ramlapi.ProcessRAML(f)
    if err != nil {
		log.Fatal(err)
	}
ramlapi.Build(api, routerFunc)
}

func routerFunc(map [string]string) {
	router.
		Methods(data["verb"]).
		Path(data["path"]).
		Handler(RouteMap[data["handler"]])
}
