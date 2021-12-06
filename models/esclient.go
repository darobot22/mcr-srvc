package models

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var initialized = false

func GetEsCLient() *elasticsearch.Client {
	var es *elasticsearch.Client
	var err error
	if !initialized {
		es, err = elasticsearch.NewDefaultClient()
		initialized = true
	}
	log.Fatal(err)
	return es
}
