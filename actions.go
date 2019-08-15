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
	json.NewEncoder(writer).Encode(result)
}

func responseMovie(writer http.ResponseWriter, status int, result Movie){
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(result)
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

func MovieUpdate(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) // recogemos todas las variables de la url
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		writer.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movie_id)
	decoder := json.NewDecoder(request.Body)

	var movie_data Movie
	err := decoder.Decode(&movie_data)

	if err != nil {
		panic(err)
		writer.WriteHeader(500)
		return
	}

	defer request.Body.Close()

	document := bson.M{"_id" : oid}
	change := bson.M{"$set" : movie_data}
	err = collection.Update(document, change)

	if err != nil {
		panic(err)
		writer.WriteHeader(404)
		return
	}

	responseMovie(writer, 404, movie_data)
}

type Message struct {
	Status string 		`json:"status"`
	Menssage string		`json:"menssage"`
}

func (this *Message) setStatus(data string)  { // *Message indica que es un puntero y no una copia de la estructura
	this.Status = data                         // para generar una copia de Message usarlo sin el "*"
}

func (this *Message) setMessage(data string)  {
	this.Menssage = data
}

func MovieDelete(writer http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request) // recogemos todas las variables de la url
	movie_id := params["id"]

	if !bson.IsObjectIdHex(movie_id) {
		writer.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(movie_id)

	err := collection.RemoveId(oid)

	if err != nil {
		writer.WriteHeader(404)
		return
	}

	message := new(Message)
	message.setStatus("success")
	message.setMessage("La pelicula con id" + movie_id + " ha sido borrada con exito")

	results := message
	writer.Header().Set("Content-Type", "application/jon")
	writer.WriteHeader(200)
	json.NewEncoder(writer).Encode(results)

}
