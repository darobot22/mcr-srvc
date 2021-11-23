package main

import (
	"encoding/json"
	"fmt"
	"log"
	"models"
	"net/http"
	"strconv"

	"github.com/rs/cors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func sendMessageToHan(rqSrvc models.RequestedService) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	if err != nil {
		fmt.Printf(err.Error())
	}

	jsonString, err := json.Marshal(rqSrvc)

	srvcString := string(jsonString)

	topic := "srvcs-topic1"
	for _, word := range []string{string(srvcString)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
	p.Flush(15 * 1000)
	return
}

func Add(w http.ResponseWriter, r *http.Request) {
	db := models.OpenConnections()
	var service models.Service
	err := json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	_, err = db.Exec("INSERT INTO services (id,name,code) VALUES ($1,$2,$3)",
		service.Id, service.Name, service.Code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	sid, err := strconv.Atoi(id)
	db := models.OpenConnections()
	db.Exec("DELETE FROM services where id = $1", sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	sid, err := strconv.Atoi(id)
	db := models.OpenConnections()
	var service models.Service
	err = json.NewDecoder(r.Body).Decode(&service)
	db.Exec("UPDATE services SET name = $1, code = $2 WHERE id = $3",
		service.Name, service.Code, sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	db := models.OpenConnections()
	var service models.Service
	sid, err := strconv.Atoi(id)
	serviceData := db.QueryRow("SELECT * FROM services WHERE id = $1", sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	err = serviceData.Scan(&service.Id, &service.Name, &service.Code)
	json.NewEncoder(w).Encode(service)
	defer db.Close()
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	db := models.OpenConnections()
	serviceData, err := db.Query("SELECT * FROM services")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	services := []models.Service{}
	var service models.Service
	for serviceData.Next() {
		err := serviceData.Scan(&service.Id, &service.Name, &service.Code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
		}
		services = append(services, service)
	}
	for _, i := range services {
		json.NewEncoder(w).Encode(i)
	}
	defer db.Close()
	defer serviceData.Close()
}

func Request(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	uid := r.FormValue("uid")

	var rqSrvc models.RequestedService
	var err error
	rqSrvc.Params = "TODO"
	rqSrvc.UserId, err = strconv.Atoi(uid)

	rqSrvc.ServiceId, err = strconv.Atoi(id)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.WriteHeader(200)
	fmt.Println(w, "StatusOK")
	sendMessageToHan(rqSrvc)
}

func main() {
	ro := mux.NewRouter().StrictSlash(true)
	ro.HandleFunc("/services", GetAll)
	ro.HandleFunc("/services/{id}", Get)
	ro.HandleFunc("/servicesadd", Add)
	ro.HandleFunc("/services/edit/{id}", Edit)
	ro.HandleFunc("/services/delete/{id}", Delete)
	ro.HandleFunc("/services/request/{id}", Request).Queries("uid", "{uid}")

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

	log.Fatal(http.ListenAndServe(":1235", handler))
}
