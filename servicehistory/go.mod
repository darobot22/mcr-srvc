module proj/servicehistory

go 1.17

require (
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.2
	github.com/rs/cors v1.8.0
    github.com/gorilla/handlers v1.5.1
    models v1.0.0
	github.com/sirupsen/logrus v1.8.1
	github.com/elastic/go-elasticsearch/v8 v8.0.0-20211123103400-5f8a17a2322f
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/confluentinc/confluent-kafka-go v1.7.0
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace models => C:/Users/User/go/src/proj/models