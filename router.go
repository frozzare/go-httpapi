// Package httpapi combine the popular httprouter package and alice to bring the best of both worlds when creating http apis.
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Param is a single URL parameter, consisting of a key and a value.
// Just a alias for httprouter.Param.
type Param = httprouter.Param

// Params is a Param-slice, as returned by the router.
// Just a alias for httprouter.Params.
type Params = httprouter.Params

// Handle is httprouter's function that can be in middlewares.
// Handle is a function that can be registered to a route to handle HTTP requests.
// Just a alias for httprouter.Handle.
type Handle = httprouter.Handle

// HandleFunc is a function that can be registered to be a route to handle HTTP requests.
type HandleFunc func(r *http.Request, ps httprouter.Params) (interface{}, interface{})

// HandleFunc2 is a function that can be registered to be a route to handle HTTP requests.
// HandleFunc2 does only have http request as a argument.
type HandleFunc2 func(r *http.Request) (interface{}, interface{})

// HandleFunc3 is a function that can be registered to be a route to handle HTTP requests.
// HandleFunc3 does only have params as a argument.
type HandleFunc3 func() (interface{}, interface{})

// HandleFunc4 is a function that can be registered to be a route to handle HTTP requests.
// HandleFunc4 does not have any arguments.
type HandleFunc4 func() (interface{}, interface{})

// ParamsFromContext pulls the URL parameters from a request context, or returns nil if none are present.
// Just a alias function for httprouter.ParamsFromContext.
func ParamsFromContext(ctx context.Context) Params {
	return httprouter.ParamsFromContext(ctx)
}

// Router represents the router.
type Router struct {
	router         *httprouter.Router
	middlewares    alice.Chain
	ResponseHandle func(HandleFunc) httprouter.Handle
}

// NewRouter creates a new router.
func NewRouter(args ...*httprouter.Router) *Router {
	r := &Router{
		middlewares: alice.New(),
	}

	if len(args) > 0 {
		r.router = args[0]
	} else {
		r.router = httprouter.New()
	}

	r.ResponseHandle = r.defaultResponseHandle

	return r
}

// Handle adds a new handle to a path and method.
func (r *Router) Handle(method, path string, handle interface{}) {
	var handler http.Handler

	// Wrap different versions of api handle functions.
	switch h := handle.(type) {
	case func(r *http.Request, ps httprouter.Params) (interface{}, interface{}):
		handler = r.wrapHandle(r.ResponseHandle(h))
	case func(r *http.Request) (interface{}, interface{}):
		handler = r.wrapHandle(r.ResponseHandle(func(r *http.Request, _ Params) (interface{}, interface{}) {
			return h(r)
		}))
	case func(ps Params) (interface{}, interface{}):
		handler = r.wrapHandle(r.ResponseHandle(func(r *http.Request, ps Params) (interface{}, interface{}) {
			return h(ps)
		}))
	case func() (interface{}, interface{}):
		handler = r.wrapHandle(r.ResponseHandle(func(r *http.Request, _ Params) (interface{}, interface{}) {
			return h()
		}))
	case func(w http.ResponseWriter, r *http.Request, ps Params):
		handler = r.wrapHandle(h)
	default:
		return
	}

	// Append middlewares using alice.
	handler = r.middlewares.Then(handler)

	// Route away!
	r.router.Handler(method, path, handler)
}

// Handler is an adapter which allows the usage of an http.Handler as a
// request handle. Just a alias function for httprouter's Handler.
func (r *Router) Handler(method, path string, handler http.Handler) {
	r.router.Handler(method, path, handler)
}

// Get is a shortcut for router.Handle("GET", path, handle).
func (r *Router) Get(path string, handle interface{}) {
	r.Handle("GET", path, handle)
}

// Head is a shortcut for router.Handle("HEAD", path, handle).
func (r *Router) Head(path string, handle interface{}) {
	r.Handle("HEAD", path, handle)
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle).
func (r *Router) Options(path string, handle interface{}) {
	r.Handle("OPTIONS", path, handle)
}

// Post is a shortcut for router.Handle("POST", path, handle).
func (r *Router) Post(path string, handle interface{}) {
	r.Handle("POST", path, handle)
}

// Put is a shortcut for router.Handle("PUT", path, handle).
func (r *Router) Put(path string, handle interface{}) {
	r.Handle("PUT", path, handle)
}

// Patch is a shortcut for router.Handle("PATCH", path, handle).
func (r *Router) Patch(path string, handle interface{}) {
	r.Handle("PATCH", path, handle)
}

// Delete is a shortcut for router.Handle("DELETE", path, handle).
func (r *Router) Delete(path string, handle interface{}) {
	r.Handle("DELETE", path, handle)
}

// Use appends a MiddlewareFunc to the chain.
// Middleware can be used to intercept or otherwise modify requests and/or responses,
// and are executed in the order that they are applied to the Router.
func (r *Router) Use(mwf ...interface{}) {
	for _, mw := range mwf {
		switch m := mw.(type) {
		case func(http.Handler) http.Handler:
			r.middlewares = r.middlewares.Append(m)
		case func(http.Handler) Handle:
			r.middlewares = r.middlewares.Append(func(h http.Handler) http.Handler {
				return r.wrapHandle(m(h))
			})
		case func(Handle) Handle:
			r.middlewares = r.middlewares.Append(func(h http.Handler) http.Handler {
				return r.wrapHandle(m(func(w http.ResponseWriter, r *http.Request, _ Params) {
					h.ServeHTTP(w, r)
				}))
			})
		case func(Handle) http.Handler:
			r.middlewares = r.middlewares.Append(func(h http.Handler) http.Handler {
				return m(func(w http.ResponseWriter, r *http.Request, _ Params) {
					h.ServeHTTP(w, r)
				})
			})
		default:
		}
	}
}

// ServeFiles serves files from the given file system root.
// Just a alias function for httprouter's ServeFiles.
// Read more: https://godoc.org/github.com/julienschmidt/httprouter#Router.ServeFiles
func (r *Router) ServeFiles(path string, root http.FileSystem) {
	r.router.ServeFiles(path, root)
}

// ServeHTTP makes the router implement the http.Handler interface.
// Just a alias function for httprouter's ServeHTTP.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) defaultResponseHandle(fn HandleFunc) Handle {
	return func(w http.ResponseWriter, r *http.Request, ps Params) {
		data, err := fn(r, ps)

		if err == nil {
			if err := json.NewEncoder(w).Encode(data); err == nil {
				return
			}
		}

		if err, ok := err.(error); ok {
			msg := err.Error()

			if msg[0] == '{' && msg[len(msg)-1] == '}' || msg[0] == '[' && msg[len(msg)-1] == ']' {
				fmt.Fprintf(w, msg)
			} else {
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}
}

// wrap wraps httprouter.Handle with http.Handler
func (r *Router) wrapHandle(next Handle) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, ParamsFromContext(r.Context()))
	})
}
