module proj/user

go 1.17

require (
	github.com/go-redis/redis/v8 v8.11.3
	github.com/gorilla/mux v1.8.0
	github.com/lib/pq v1.10.2
	github.com/rs/cors v1.8.0
	github.com/gorilla/handlers v1.5.1
	models v1.0.0
	github.com/sirupsen/logrus v1.8.1
)

require (
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace models => C:/Users/User/go/src/proj/models