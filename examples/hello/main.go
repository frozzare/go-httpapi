package main

import (
	"net/http"

	"github.com/frozzare/go-httpapi"
)

func main() {
	router := httpapi.NewRouter()

	router.Get("/hello/:name", func(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
		return map[string]string{
			"hello": ps.ByName("name"),
		}, nil
	})

	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}
