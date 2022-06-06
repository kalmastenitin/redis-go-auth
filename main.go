package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var CONN *redis.Client

func loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	username := r.Header.Get("username")
	password := r.Header.Get("password")

	newUUID, _ := uuid.NewUUID()

	if username == "admin" && password == "admin" {
		log.Println("credentials are correct")
		_, err := CONN.Do(ctx, "SET", newUUID.String(), "encrypted_secret_data", "ex", 30).Result()
		if err != nil {
			log.Fatal(err)
		}
		response, _ := json.Marshal(map[string]string{"value": newUUID.String()})
		w.Write(response)
		return
	}
}

func tokenValidator(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	token := r.Header.Get("token")
	val, err := CONN.Get(ctx, token).Result()
	log.Println(err)
	if val != "" {
		response, _ := json.Marshal(map[string]string{"response": "success"})
		w.Write(response)
		return
	}
	w.WriteHeader(http.StatusUnauthorized)
	response, _ := json.Marshal(map[string]string{"response": "unauthorized"})
	w.Write(response)
	return
}

func main() {
	// ctx := context.Background()
	CONN = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "sec_rEtPass", // no password set
		DB:       0,             // use default DB
	})

	r := mux.NewRouter()
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/validate", tokenValidator).Methods("GET")
	log.Println("starting server ...")
	http.ListenAndServe(":8000", r)
}
