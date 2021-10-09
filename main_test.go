package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// test for main function
func TestMain(t *testing.T) {
}

// test to check hashing
func TestGetHashed256(t *testing.T) {
	expectedString := "008c70392e3abfbd0fa47bbc2ed96aa99bd49e159727fcba0f2e6abeb3a9d601"
	//computed hash for "Password123"
	actualString := getHashed256("Password123")
	if actualString != expectedString {
		t.Errorf("Expected String(%s) is not same as"+
			" actual string (%s)", expectedString, actualString)
	}
}

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/posts", GetPostByid).Methods("GET")
	return router
}

// test to check get request of post
func TestGetPost(t *testing.T) {
	request, _ := http.NewRequest("GET", "/posts1", nil)
	response := httptest.NewRecorder()
	Router().ServeHTTP(response, request)
	status := response.Code
	if status == 200 {
		t.Errorf("Expected code not same as %v", http.StatusOK)
	}
}
