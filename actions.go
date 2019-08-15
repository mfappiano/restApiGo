package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)
var collection = getSession().DB("siacDB").C("movies")

func getSession() *mgo.Session{
	session, err := mgo.Dial("mongodb://localhost:27017")

	if (err != nil ){
		panic(err)
	}

	return session
}

func responseMovies(writer http.ResponseWriter, status int, result []Movie){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(results)
}

func responseMovie(writer http.ResponseWriter, status int, result Movie){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(results)
}

func Index(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hola mundo desde mi server go")
}

func MovieList(writer http.ResponseWriter, request *http.Request) {
	var results []Movie
	err := collection.Find(nil).Sort("") .All(&results)

	if err != nil {
		log.Fatal(err)
		panic(err)
	} else {
		fmt.Println("Resultados ", request)
	}
	responseMovies(writer, 200, results)
}

func MovieShow(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) // recogemos todas las variables de la url
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		writer.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movie_id)

	results := Movie{}

	err := collection.FindId(oid).One(&results)

	if err != nil {
		writer.WriteHeader(404)
		return
	} else {
		responseMovie(writer, 404, results)

	}
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
	responseMovie(writer, 404, movie_data)
}