package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"models"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var client = redis.NewClient(&redis.Options{
	Addr:     "host.docker.internal:6379",
	Password: "",
	DB:       0,
})

func GetRedis() redis.Client {
	return *client
}

func Add(w http.ResponseWriter, r *http.Request) {
	db := models.OpenConnections()
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = db.Exec("INSERT INTO users (id,name,password,email,phone) VALUES ($1,$2,$3,$4,$5)",
		user.Id, user.Name, user.Password, user.Email, user.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var ctx = context.Background()
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	client := GetRedis()
	db := models.OpenConnections()
	var user models.User
	userJson, err := client.Get(ctx, id).Result()
	if err == redis.Nil {
		fmt.Println("Redis not stored")
		uid, err2 := strconv.Atoi(id)
		userData := db.QueryRow("SELECT * FROM users WHERE id = $1", uid)
		err2 = userData.Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Phone)
		if err2 != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err2 = json.NewEncoder(w).Encode(user)
		userResult, err2 := json.Marshal(user)
		err2 = client.Set(ctx, id, userResult, 300*time.Second).Err()
		if err2 != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
		}
		defer db.Close()
	} else {
		fmt.Println("Redis stored")
		targets := models.User{}

		err := json.Unmarshal([]byte(userJson), &targets)
		err = json.NewEncoder(w).Encode(targets)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	db := models.OpenConnections()
	usersData, err := db.Query("SELECT * FROM users")
	users := []models.User{}
	var user models.User
	for usersData.Next() {
		err = usersData.Scan(&user.Id, &user.Name, &user.Password, &user.Email, &user.Phone)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		users = append(users, user)
	}
	for _, i := range users {
		err = json.NewEncoder(w).Encode(i)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
	}
	defer db.Close()
	defer usersData.Close()
}

func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	db := models.OpenConnections()
	uid, err := strconv.Atoi(id)
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	_, err = db.Exec("UPDATE users SET name = $1, password = $2, email = $3, phone = $4 where id = $5",
		user.Name, user.Password, user.Email, user.Phone, uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	db := models.OpenConnections()
	uid, err := strconv.Atoi(id)
	_, err = db.Exec("DELETE FROM users where id = $1", uid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer db.Close()
}

func main() {
	ro := mux.NewRouter().StrictSlash(true)
	ro.HandleFunc("/users", GetAll)
	ro.HandleFunc("/users/{id}", Get)
	ro.HandleFunc("/usersadd", Add)
	ro.HandleFunc("/users/edit/{id}", Edit)
	ro.HandleFunc("/users/delete/{id}", Delete)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
	})

	handler := c.Handler(ro)
	log.SetFormatter(&log.JSONFormatter{})
	f, err := os.OpenFile("servicehistory.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Error(err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Fatal(http.ListenAndServe(":1234", handler))
}
