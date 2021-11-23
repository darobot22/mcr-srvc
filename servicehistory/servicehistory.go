package main

import (
	"encoding/json"
	"fmt"
	"log"
	"models"
	"net/http"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func HistoryGetAll(w http.ResponseWriter, r *http.Request) {
	db := models.OpenConnections()
	serviceData, err := db.Query("SELECT * FROM servicehistory")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	services := []models.ServiceHistory{}
	var service models.ServiceHistory
	for serviceData.Next() {
		err := serviceData.Scan(&service.Id, &service.ServiceCode, &service.ServiceName, &service.UserId,
			&service.UserName, &service.CreateDate, &service.ResultData, &service.ExecutionDate)
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

func HistoryEdit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	sid, err := strconv.Atoi(id)
	db := models.OpenConnections()
	var servicehistory models.ServiceHistory
	err = json.NewDecoder(r.Body).Decode(&servicehistory)
	db.Exec("UPDATE servicehistory SET resultdata = $1, executiondate = $2 WHERE id = $3",
		servicehistory.ResultData, servicehistory.ExecutionDate, sid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func HistoryGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
	}
	sid, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	db := models.OpenConnections()
	serviceData, err := db.Query("SELECT * FROM servicehistory WHERE userid = $1", sid)
	services := []models.ServiceHistory{}
	var service models.ServiceHistory
	for serviceData.Next() {
		err := serviceData.Scan(&service.Id, &service.ServiceCode, &service.ServiceName, &service.UserId,
			&service.UserName, &service.CreateDate, &service.ResultData, &service.ExecutionDate)

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

func saveService(srvc string) {
	var rqsrvcs models.RequestedService
	err := json.Unmarshal([]byte(srvc), &rqsrvcs)
	var historyService models.ServiceHistory
	historyService.UserId = rqsrvcs.UserId
	db := models.OpenConnections()
	serviceData := db.QueryRow("SELECT name,code FROM services WHERE id = $1", rqsrvcs.ServiceId)
	if err != nil {
		return
	}
	err = serviceData.Scan(&historyService.ServiceName, &historyService.ServiceCode)
	serviceData = db.QueryRow("SELECT name FROM users WHERE id = $1", rqsrvcs.UserId)
	if err != nil {
		return
	}
	err = serviceData.Scan(&historyService.UserName)

	historyService.CreateDate = time.Now().Format("01-02-2006")
	historyService.ResultData = rqsrvcs.Params
	historyService.ExecutionDate = "NOT SET"

	_, err = db.Exec("INSERT INTO servicehistory (servicecode,servicename,userid,username,createdate,resultdata) VALUES ($1,$2,$3,$4,$5,$6)",
		historyService.ServiceCode, historyService.ServiceName,
		historyService.UserId, historyService.UserName, historyService.CreateDate, historyService.ResultData)
	if err != nil {
		panic(err)
	}
}

func HandleHistory() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"srvcs-topic2"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			srvc := string(msg.Value)
			saveService(srvc)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func main() {
	ro := mux.NewRouter().StrictSlash(true)
	ro.HandleFunc("/history/", HistoryGetAll)
	ro.HandleFunc("/history/{id}", HistoryGetUser)
	ro.HandleFunc("/history/edit/{id}", HistoryEdit)

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

	go HandleHistory()
	log.Fatal(http.ListenAndServe(":1236", handler))
}
