package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ServiceHistory struct {
	Id            int
	ServiceCode   string
	ServiceName   string
	UserId        int
	UserName      string
	CreateDate    string
	ResultData    string
	ExecutionDate string
}

type RequestedService struct {
	UserId    int
	ServiceId int
	Params    string
}

func OpenConnections() *sql.DB {
	connStr := "host=host.docker.internal port=5432 user=postgres password=fkubifkom10 dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	return db
}

func HistoryGetAll(w http.ResponseWriter, r *http.Request) {
	db := OpenConnections()
	serviceData, err := db.Query("SELECT * FROM servicehistory")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
	}
	services := []ServiceHistory{}
	var service ServiceHistory
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
	db := OpenConnections()
	var servicehistory ServiceHistory
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
	db := OpenConnections()
	serviceData, err := db.Query("SELECT * FROM servicehistory WHERE userid = $1", sid)
	services := []ServiceHistory{}
	var service ServiceHistory
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

func saveService(srvc string){
	var rqsrvcs RequestedService
	err := json.Unmarshal([]byte(srvc), &rqsrvcs)
	var historyService ServiceHistory
	historyService.UserId = rqsrvcs.UserId
	db := OpenConnections()
	serviceData := db.QueryRow("SELECT name,code FROM services WHERE id = $1", rqsrvcs.ServiceId)
	if err != nil {
		panic(err)
	}
	err = serviceData.Scan(&historyService.ServiceName, &historyService.ServiceCode)
	serviceData = db.QueryRow("SELECT name FROM users WHERE id = $1", rqsrvcs.UserId)
	if err != nil {
		panic(err)
	}
	err = serviceData.Scan(&historyService.UserName)

	historyService.CreateDate = time.Now().Format("01-02-2006")
	historyService.ResultData = rqsrvcs.Params
	historyService.ExecutionDate = "NOT SET"

	_, err = db.Exec("INSERT INTO servicehistory (servicecode,servicename,userid,username,createdate,resultdata) VALUES ($1,$2,$3,$4,$5,$6)",
		historyService.ServiceCode, historyService.ServiceName,
		historyService.UserId, historyService.UserName, historyService.CreateDate, historyService.ResultData, )
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

func main(){
	ro := mux.NewRouter().StrictSlash(true)
	ro.HandleFunc("/history/", HistoryGetAll)
	ro.HandleFunc("/history/{id}", HistoryGetUser)
	ro.HandleFunc("/history/edit/{id}", HistoryEdit)
	go HandleHistory()
	log.Fatal(http.ListenAndServe(":1236", ro))
}
