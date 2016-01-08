package main

const handlerHead = `package main

import (
	"encoding/json"
	"net/http"
)
`

const mapStart = `
var RouteMap = map[string]http.HandlerFunc{
`

const mapEntry = `
	"{{.Name}}":         {{.Struct}},
`

const mapEnd = `
}
`

const handlerText = `
// {{.Name}} - handler for URI {{.Path}} HTTP verb {{.Verb}}
// {{.Doc}}
func {{.Name}}(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(map[string]string{
		"message": "{{.Name}}{{.Verb}}",
 	})
 	w.Write(json)
}
`
