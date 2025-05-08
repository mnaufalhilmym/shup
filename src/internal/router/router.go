package router

import (
	"net/http"

	"shup.hilmy.dev/src/internal/controller"
)

type Router struct {
	mux            *http.ServeMux
	fileController *controller.File
}

func New(mux *http.ServeMux, fileController *controller.File) *Router {
	return &Router{mux, fileController}
}

func (c *Router) Configure() http.Handler {
	c.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			c.fileController.Save(w, r)
		case http.MethodGet:
			c.fileController.Get(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return c.mux
}
