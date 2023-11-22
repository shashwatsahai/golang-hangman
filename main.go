package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/shashwatsahai/golang-hangman/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	db, err := utils.MongoClient()
	if err != nil {
		panic(err)
	}
	hangmancol := db.Collection("hangman")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	length, err := hangmancol.CountDocuments(ctx, bson.M{})
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error counting documents", http.StatusInternalServerError)
		return
	}
	// rand.Seed(time.Now().UnixNano())
	selected := rand.Intn(int(length))

	fmt.Println("GET request for /word", selected)

	findOptions := options.Find().SetSkip(int64(selected)).SetLimit(1)

	// Define a filter (empty in this example to match all documents)
	filter := bson.D{{}}

	// Perform the find operation
	cursor, err := hangmancol.Find(context.TODO(), filter, findOptions)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error counting documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())
	var result bson.M
	if cursor.Next(context.TODO()) {
		err := cursor.Decode(&result)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error counting documents", http.StatusInternalServerError)
			return
		}
		fmt.Println("result", result["word"])
	}

	fmt.Fprint(w, result["word"])

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
		document := bson.D{{Key: "word", Value: words[i]}}
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
