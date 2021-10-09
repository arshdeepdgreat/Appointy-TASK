package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Structure of User
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

// mongo client
var client *mongo.Client

// Function to hash passcode so that it cant be reverse engineered
//return string
func getHashed256(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", hash)
}

//Create user endpoint notes
// 1. routes for /users path
// 2. accepts json data via post request
// 3. encodes the passcode
// 4. inserts 1 entry into database
func Createuser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(request.Body).Decode(&user)
	var Newpass = user.Password
	user.Password = getHashed256(Newpass)
	collection := client.Database("appointytask").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(response).Encode(result)
}

//GetByid user endpoint notes
// 1. routes for /users/{id} path
// 2. accepts json data via get request
// 3. Searches ID in the database (findone)
func GetByid(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var user User
	collection := client.Database("appointytask").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, User{ID: id}).Decode(&user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(user)
}

// Structure of Post
type Post struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption   string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Imagepath string             `json:"imagepath,omitempty" bson:"imagepath,omitempty"`
	Timestamp primitive.DateTime `json:"_timestamp,omitempty" bson:"_timestamp,omitempty"`
	User      string             `json:"user,omitempty" bson:"user,omitempty"`
}

//GetPostByid endpoint notes
// 1. routes for /posts/{id} path
// 2. accepts json data via get request
// 3. Searches ID in the database (findone)
func GetPostByid(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var post Post
	collection := client.Database("appointytask").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, Post{ID: id}).Decode(&post)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(post)
}

//GeAlltPostByUser endpoint notes
// 1. routes for /posts/users/{id} path
// 2. accepts json data via get request
// 3. Searches ID in the database (findall)
// 4. output is limitied (Pagination) to 10 entries only
func Getallpostbyuser(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var posts []Post
	params := mux.Vars(request)
	var ids = params["id"]

	collection := client.Database("appointytask").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := collection.Find(ctx, bson.M{"user": ids})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
		//len reaches 10 so break
		if len(posts) == 10 {
			break
		}
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}
	json.NewEncoder(response).Encode(posts)
}

//Create Post endpoint notes
// 1. routes for /posts path
// 2. accepts json data via post request
// 3. inserts 1 entry into database

func Createpost(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(request.Body).Decode(&post)

	post.Timestamp = primitive.DateTime(time.Now().Unix())
	//posting timestamp is put exactly when data is stored

	collection := client.Database("appointytask").Collection("posts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(response).Encode(result)
}

func main() {
	fmt.Println("App started...")
	//App started notification

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	router := mux.NewRouter()

	// USER Endpoints
	router.HandleFunc("/users", Createuser).Methods("POST")
	router.HandleFunc("/users/{id}", GetByid).Methods("GET")

	// POST Endpoints
	router.HandleFunc("/posts", Createpost).Methods("POST")
	router.HandleFunc("/posts/{id}", GetPostByid).Methods("GET")

	//FIND POSTS BY USER
	router.HandleFunc("/posts/users/{id}", Getallpostbyuser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
	//Sample data for post
	// {
	// 	"Caption":"Sample Caption",
	// 	"Imagepath":"C:\\images\\pic.jpeg",
	// 	"User":"61613fa0d11cdcdc435c511f"
	// }

	//Sample Data for user
	// {
	// 	"Name":"Arshdeep",
	// 	"Email":"arshdeepdgreat@gmail.com"
	// 	"Password":"Password123"
	// }

}
