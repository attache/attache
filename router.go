package attache

import "github.com/go-chi/chi"

type Router interface {
	chi.Router
}

type router struct {
	chi.Mux
}
