package httpapi

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestRouter(t *testing.T) {
	router := NewRouter()

	router.ResponseHandle = func(fn HandleFunc) Handle {
		return func(w http.ResponseWriter, r *http.Request, ps Params) {
			data, out := fn(r, ps)

			if err, ok := out.(error); ok {
				t.Fatal(err)
			}

			if "Hello gopher" != data.(string) {
				t.Fatal("Wrong result")
			}
		}
	}

	routed := false
	router.Handle("GET", "/user/:name", func(r *http.Request, ps Params) (interface{}, interface{}) {
		routed = true
		want := Params{Param{Key: "name", Value: "gopher"}}
		if !reflect.DeepEqual(ps, want) {
			return nil, fmt.Errorf("Wrong wildcard values: want %v, got %v", want, ps)
		}

		return fmt.Sprintf("Hello %s", want.ByName("name")), nil
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	router.ServeHTTP(w, req)

	if !routed {
		t.Fatal("Routing failed")
	}
}

func TestRouterAPI(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	router := NewRouter()
	router.Get("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		get = true
		return nil, nil
	})
	router.Head("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		head = true
		return nil, nil
	})
	router.Options("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		options = true
		return nil, nil
	})
	router.Post("/POST", func(r *http.Request, ps Params) (interface{}, interface{}) {
		post = true
		return nil, nil
	})
	router.Put("/PUT", func(r *http.Request, ps Params) (interface{}, interface{}) {
		put = true
		return nil, nil
	})
	router.Patch("/PATCH", func(r *http.Request, ps Params) (interface{}, interface{}) {
		patch = true
		return nil, nil
	})
	router.Delete("/DELETE", func(r *http.Request, ps Params) (interface{}, interface{}) {
		delete = true
		return nil, nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}

func TestHandleFunc2(t *testing.T) {
	var get bool

	router := NewRouter()
	router.Get("/GET", func(r *http.Request) (interface{}, interface{}) {
		get = true
		return nil, nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}
}

func TestHandleFunc3(t *testing.T) {
	var get bool

	router := NewRouter()
	router.Get("/GET/:name", func(ps Params) (interface{}, interface{}) {
		get = ps.ByName("name") == "fredrik"
		return nil, nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET/fredrik", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}
}

func TestHandleFunc4(t *testing.T) {
	var get bool

	router := NewRouter()
	router.Get("/GET", func() (interface{}, interface{}) {
		get = true
		return nil, nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}
}

func TestHandleFunc5(t *testing.T) {
	var get bool

	router := NewRouter()
	router.Get("/GET", func(w http.ResponseWriter, r *http.Request, _ Params) {
		get = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}
}

func TestHandleFunc6(t *testing.T) {
	var get bool
	var get2 bool

	router := NewRouter()
	router.Get("/GET", func(w http.ResponseWriter, r *http.Request) {
		get = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	router.Get("/NAME/:name", func(w http.ResponseWriter, r *http.Request) {
		ps := ParamsFromContext(r.Context())

		if ps.ByName("name") == "fredrik" {
			get2 = true
		}
	})

	r, _ = http.NewRequest("GET", "/NAME/fredrik", nil)
	router.ServeHTTP(w, r)
	if !get2 {
		t.Error("routing GET with Params failed")
	}
}

func TestGroup(t *testing.T) {
	var get bool
	var get2 bool

	router := NewRouter()
	router.Get("/GET", func(w http.ResponseWriter, r *http.Request, _ Params) {
		get = true
	})

	router.Group("/FOO").Group("/BAR").Get("/GET", func(w http.ResponseWriter, r *http.Request, _ Params) {
		get2 = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("GET", "/FOO/BAR/GET", nil)
	router.ServeHTTP(w, r)
	if !get2 {
		t.Error("routing GET GROUP failed")
	}
}

func TestRouterAPIGroup(t *testing.T) {
	var get, head, options, post, put, patch, delete bool

	base := NewRouter()

	router := base.Group("/FOO")

	router.Get("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		get = true
		return nil, nil
	})
	router.Head("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		head = true
		return nil, nil
	})
	router.Options("/GET", func(r *http.Request, ps Params) (interface{}, interface{}) {
		options = true
		return nil, nil
	})
	router.Post("/POST", func(r *http.Request, ps Params) (interface{}, interface{}) {
		post = true
		return nil, nil
	})
	router.Put("/PUT", func(r *http.Request, ps Params) (interface{}, interface{}) {
		put = true
		return nil, nil
	})
	router.Patch("/PATCH", func(r *http.Request, ps Params) (interface{}, interface{}) {
		patch = true
		return nil, nil
	})
	router.Delete("/DELETE", func(r *http.Request, ps Params) (interface{}, interface{}) {
		delete = true
		return nil, nil
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/FOO/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("HEAD", "/FOO/GET", nil)
	router.ServeHTTP(w, r)
	if !head {
		t.Error("routing HEAD failed")
	}

	r, _ = http.NewRequest("OPTIONS", "/FOO/GET", nil)
	router.ServeHTTP(w, r)
	if !options {
		t.Error("routing OPTIONS failed")
	}

	r, _ = http.NewRequest("POST", "/FOO/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/FOO/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/FOO/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/FOO/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}
}
