This package is designed to work with [RAML](http://raml.org) files and build the appropriate handlers and routing rules in Go applications.

[![GoDoc](https://godoc.org/github.com/EconomistDigitalSolutions/ramlapi?status.svg)](https://godoc.org/github.com/EconomistDigitalSolutions/ramlapi)

The ramlapi codebase contains two packages:

* Ramlapi - used to parse a RAML file and wire it up to a router.
* Ramlgen - used to parse a RAML file and write a set of HTTP handlers.

#### RAML Compatibility

The current version of the Ramlapi and Ramlgen packages supports *most* of the 0.8 RAML specification.

Our intention is to implement further support as:

1. Preliminary 1.0 support
2. Full 1.0 support
3. Additional 0.8 support

Enhancing 0.8 support is a low priority as users are strongly urged to migrate to 1.0 as soon as possible

#### HOW TO RAML-GEN

1. Build your API design in RAML.
2. Inside your Go project, run `raml-gen --ramlfile=<file>`.
3. Copy the resulting `handlers_gen.go` file to the correct location.

You now have a set of HTTP handlers built from your RAML specification.

Right now raml-gen only supports routers that use the standard http.Handlerfunc handlers. If you want to use something like
Echo or HttpRouter you'll need to make some amendments as described in the examples below.


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
	path := ep.Path

	for _, up := range ep.URIParameters {
		if up.Pattern != "" {
			path = strings.Replace(
				path,
				fmt.Sprintf("{%s}", up.Key),
				fmt.Sprintf("{%s:%s}", up.Key, up.Pattern),
				1)
		}
	}

	route := router.
		Methods(ep.Verb).
		Path(path).
		Handler(RouteMap[ep.Handler])

	for _, qp := range ep.QueryParameters {
		if qp.Required {
			if qp.Pattern != "" {
				route.Queries(qp.Key, fmt.Sprintf("{%s:%s}", qp.Key, qp.Pattern))
			} else {
				route.Queries(qp.Key, "")
			}
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
