package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shashwatsahai/golang-hangman/utils"
	"go.mongodb.org/mongo-driver/bson"
)

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	params, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		panic(err)
// 	}
// 	var data map[string]interface{}
// 	paramvalue := json.Unmarshal(params, &data)
// 	s, t, d := reflect.TypeOf(params), reflect.TypeOf(paramvalue), reflect.TypeOf(data)
// 	fmt.Println("Param", paramvalue, s, t, d, data["name"])
// 	fmt.Fprint(w, "Hello, this is the home page!")
// }

//	func loggingMiddleware(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			fmt.Print("New Request")
//			next.ServeHTTP(w, r)
//		})
//	}
func wordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	// Handle GET request logic here
	fmt.Fprint(w, "GET request for /word")

}

func handlePost(w http.ResponseWriter, r *http.Request) {
	// Handle POST request logic here
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var data map[string][]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	db, err := utils.MongoClient()
	if err != nil {
		panic(err)
	}
	hangmancol := db.Collection("hangman")
	fmt.Printf("Received POST request with data: %v\n", data)
	words := data["words"]
	// s := make([]interface{}, len(words))
	for i := 0; i < len(words); i++ {
		fmt.Println(words[i])
		// s[i] = words[i]
		document := bson.D{{"word", words[i]}}
		result, err := hangmancol.InsertOne(context.TODO(), document)
		if err != nil {
			panic(err)
		}
		println(result)

	}
	// hangmancol.InsertMany(context.TODO(), s)
	fmt.Fprint(w, "POST request for /word")
}
func main() {
	server := http.NewServeMux()
	// server.Handle("/", loggingMiddleware(http.HandlerFunc(homeHandler)))
	server.HandleFunc("/word", wordHandler) //ger and post word)
	http.ListenAndServe(":3000", server)
}
