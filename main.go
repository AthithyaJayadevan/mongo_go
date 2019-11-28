package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Person struct definition
type Person struct {
	//ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Firstname string `json:"firstname,omitempty" bson:"firstname, omitempty"`
	Lastname  string `json:"lastname,omitempty" bson:"lastname, omitempty"`
}

var client *mongo.Client

//Createperson func to create DB entry through http
func Createperson(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)

	fmt.Print(person)

	collection := client.Database("test1").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func Getperson(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("content-type", "application/json")
	var person Person
	params := mux.Vars(r)
	id, _ := params["firstname"]
	collections := client.Database("test1").Collection("customers")
	ctx, _ := context.WithTimeout(context.Background(), 50*time.Second)
	err := collections.FindOne(ctx, bson.M{"firstname": id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(person)

}

func main() {
	fmt.Println("Starting application...")
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	//fmt.Printf("Connecting to DB\n")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, _ = mongo.Connect(ctx, clientOptions)

	router := mux.NewRouter()
	fmt.Printf("Connecting to DB\n")

	router.HandleFunc("/person", Createperson).Methods("POST")
	router.HandleFunc("/person/{firstname}", Getperson).Methods("GET")

	http.ListenAndServe(":8000", router)
}
