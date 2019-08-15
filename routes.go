package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	name string
	method string
	pattern string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true) // las urls son amigables usando /
	for _, route := range routes {
		router.
			Name(route.name).
			Methods(route.method).
			Path(route.pattern).
			Handler(route.HandlerFunc)
	}
	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"MovieList",
		"GET",
		"/peliculas",
		MovieList,
	},
	Route{
		"MovieShow",
		"GET",
		"/pelicula/{id}",
		MovieShow,
	},
	Route{
		name:        "MoviedAdd",
		method:      "Post",
		pattern:     "/pelicula",
		HandlerFunc: MovieAdd,
	},
	Route{
		name:        "MovieUpdate",
		method:      "PUT",
		pattern:     "/pelicula/{id}",
		HandlerFunc: MovieUpdate,
	},
	Route{
		name:        "MovieDelete",
		method:      "DELETE",
		pattern:     "/pelicula/{id}",
		HandlerFunc: MovieDelete,
	},
}
