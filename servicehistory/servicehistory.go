package main

import (
	"context"
	"encoding/json"
	"fmt"
	"models"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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
		log.Error(err)
	}
	services := []models.ServiceHistory{}
	var service models.ServiceHistory
	for serviceData.Next() {
		err := serviceData.Scan(&service.Id, &service.ServiceCode, &service.ServiceName, &service.UserId,
			&service.UserName, &service.CreateDate, &service.ResultData, &service.ExecutionDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			log.Error(err)
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
		log.Error(err)
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
		log.Error(err)
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
			log.Error(err)
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
			log.Warn("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func HistorySearch(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("query")
	searchRequest := fmt.Sprintf(`{"query": {"multi_match": {"query": "%s", "fields": 
	["id", "servicecode", "servicename", "userid","username","createdate","reultdata",executiondate]}}}`,
		query)
	es := models.GetEsCLient()
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("servicehistory"),
		es.Search.WithBody(strings.NewReader(searchRequest)),
		es.Search.WithPretty(),
	)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(res)
}

func main() {
	ro := mux.NewRouter().StrictSlash(true)
	ro.HandleFunc("/history/", HistoryGetAll)
	ro.HandleFunc("/history/{id}", HistoryGetUser)
	ro.HandleFunc("/history/edit/{id}", HistoryEdit)
	ro.HandleFunc("/history/search", HistorySearch).Queries("query", "{query}")

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

	go HandleHistory()
	log.Fatal(http.ListenAndServe(":1236", handler)).Info("Server started")
}
