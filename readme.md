# httpapi [![Build Status](https://travis-ci.org/frozzare/go-httpapi.svg?branch=master)](https://travis-ci.org/frozzare/go-httpapi) [![GoDoc](https://godoc.org/github.com/frozzare/go-httpapi?status.svg)](http://godoc.org/github.com/frozzare/go-httpapi) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/go-httpapi)](https://goreportcard.com/report/github.com/frozzare/go-httpapi)

> Work in progress

A http router for building http apis in Go based on [httprouter](https://github.com/julienschmidt/httprouter) and [alice](https://github.com/justinas/alice).

The router works almost the same way as [httprouter](https://github.com/julienschmidt/httprouter) does with some changes:

- Uppercase method functions is capitalized instead.
- `HandleFunc` has two arguments instead of three (request and params).
- All `HandleFunc` returns the data and/or a error that the response handle will handle.
- `HandleFunc2` has one argument instead of three (request).
- `HandleFunc3` has zero arguments instead of three.
- Default response handler that response with JSON. Can be replaced by a custom handler function.
- Not all methods exists on `httpapi.Router` struct as `httprouter.Router` has, e.g `HandlerFunc` does not exist.
- Better support for middlewares with [alice](https://github.com/justinas/alice).

## Installation

```
go get -u github.com/frozzare/go-httpapi
```

## Usage

Example code:

```go
router := httpapi.NewRouter()

router.Get("/hello/:name", func(r *http.Request, ps httpapi.Params) (interface{}, interface{}) {
    return map[string]string{
        "hello": ps.ByName("name"),
    }, nil
})

http.Handle("/", router)
http.ListenAndServe(":3000", nil)
```

Example response:

```json
GET /hello/fredrik
{
    "hello": "fredrik"
}
```

To configure [httprouter](https://github.com/julienschmidt/httprouter) you just pass it as argument to `NewRouter`:

```go
router := httpapi.NewRouter(&httprouter.Router{
    RedirectTrailingSlash: true,
})
```

To modify the response handle that takes in `HandleFunc`, `HandleFunc2` and `HandleFunc3` is wrapped with `HandleFunc`:

```go
router := httpapi.NewRouter()

router.ResponseHandle = func(fn httpapi.HandleFunc) httpapi.Handle {
    return func(w http.ResponseWriter, r *http.Request, ps httpapi.Params) {
        data, out := fn(r, ps)

        // and so on...
    }
}
```

Both return values are returned as interfaces to support more than just than the error type.

## Middlewares

```go
router := httpapi.NewRouter()

// with standard http handler.
router.Use(func(h http.Handler) http.Handler {
    fmt.Println("Hello, world")
    return h
})

// with httprouter's handle.
router.Use(func(h httpapi.Handle) httpapi.Handle {
    fmt.Println("Hello, world")
    return h
})
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)