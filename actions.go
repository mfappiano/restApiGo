package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"net/http"
)
var movies = Movies{}
/*
var movies = Movies{
		Movie{"sin limites", 2233, "No se"},
		Movie{"no tengo imaginacion", 2334, "pepe"},
		Movie{"dos mas dos cinco", 1245, "lalala"},
	}*/
var collection = getSession().DB("siacDB").C("movies")
func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://localhost:27017")

	if (err != nil ){
		panic(err)
	}

	return session
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hola mundo desde mi server go")
}

func MovieList(writer http.ResponseWriter, request *http.Request) {
	json.NewEncoder(writer).Encode(movies)
}

func MovieShow(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) // recogemos todas las variables de la url
	movie_id := params["id"]
	fmt.Fprintf(writer, "Has cargado la pelicula numero %s", movie_id)
}

func MovieAdd(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var movie_data Movie
	err := decoder.Decode(&movie_data)

	if(err != nil ){
		panic(err)
	}

	defer request.Body.Close() // cierro y limpio body

	err2 := collection.Insert(movie_data)

	if(err2 != nil ){
		writer.WriteHeader(500)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(movie_data)
}