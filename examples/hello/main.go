package main

import (
	"fmt"
	"net/http"

	"github.com/frozzare/go-httpapi"
)

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "group router")
		next.ServeHTTP(w, r)
	})
}

func main() {
	router := httpapi.NewRouter()

	router.Get("/hello/:name", func(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
		return map[string]string{
			"hello": ps.ByName("name"),
		}, nil
	})

	a := router.Group("/foo")
	a.Use(middleware)
	a.Get("/hello/:name", func(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
		return map[string]string{
			"hello": ps.ByName("name"),
		}, nil
	})

	http.Handle("/", router)
	http.ListenAndServe(":3000", nil)
}
