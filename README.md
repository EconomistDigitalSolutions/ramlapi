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

##### STANDARD LIBRARY

```go

var RouteMap = map[string]http.HandlerFunc{

  "Root":    Root,
  "Version": Version,
}

func main() {
  router := http.NewServeMux()
  api, _ := ramlapi.Process("api.raml")

  ramlapi.Build(api, routerFunc)
  log.Fatal(http.ListenAndServe(":9494", router))
}

func routerFunc(ep *ramlapi.Endpoint) {
  handler := http.HandlerFunc(RouteMap[ep.Handler])
  router.Handle(ep.Path, handler)
}

```

##### PAT

```go

var RouteMap = map[string]http.HandlerFunc{

	"Root":    Root,
	"Version": Version,
}

func Version(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func Root(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func main() {
  router := pat.New()
  api, _ := ramlapi.Process("api.raml")

  ramlapi.Build(api, routerFunc)
  log.Fatal(http.ListenAndServe(":9494", router))
}

func routerFunc(ep *ramlapi.Endpoint) {
	router.Add(ep.Verb, ep.Path, RouteMap[ep.Handler])
}
```

##### GORILLA MUX

```go

var RouteMap = map[string]http.HandlerFunc{

	"Root":    Root,
	"Version": Version,
}

func Version(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func Root(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func main() {
  router := mux.NewRouter().StrictSlash(true)
  api, _ := ramlapi.Process("api.raml")

  ramlapi.Build(api, routerFunc)
  log.Fatal(http.ListenAndServe(":9494", router))

}

func routerFunc(ep *ramlapi.Endpoint) {
	route := router.
		Methods(ep.Verb).
		Path(ep.Path).
		Handler(RouteMap[ep.Handler])

	for _, param := range ep.QueryParameters {
		if param.Pattern != "" {
			route.Queries(param.Key, fmt.Sprintf("{%s:%s}", param.Key, param.Pattern))
		} else {
			route.Queries(param.Key, "")
		}
	}
}
```

##### ECHO

```go

var RouteMap = map[string]func(c *echo.Context) error{

	"Root":    Root,
	"Version": Version,
}

func Version(c *echo.Context) error {
  return c.String(http.StatusOK, "VERSION")
}

func Root(c *echo.Context) error {
  return c.String(http.StatusOK, "HOME")
}

func main() {
	router := echo.New()

  api, _ := ramlapi.ProcessRAML("api.raml")

	ramlapi.Build(api, routerFunc)

	router.Run(":9494")
}

func routerFunc(ep *ramlapi.Endpoint) {
	switch ep.Verb {
	case "GET":
		router.Get(ep.Path, RouteMap[ep.Handler])
	}
}
```

##### HTTPROUTER

```go

var RouteMap = map[string]httprouter.Handle{

	"Root":    Root,
	"Version": Version,
}

func Version(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "VERSION\n")
}

func Root(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  fmt.Fprint(w, "HOME\n")
}

func main() {
	api, _ := ramlapi.ProcessRAML("api.raml")

	router := httprouter.New()
	ramlapi.Build(api, routerFunc)

	log.Fatal(http.ListenAndServe(":9494", router))
}

func routerFunc(ep *ramlapi.Endpoint) {
	switch ep.Verb {
	case "GET":
		router.GET(ep.Path, RouteMap[ep.Handler])
	}
}
```