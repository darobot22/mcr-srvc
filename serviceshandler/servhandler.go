package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type RequestedService struct {
	UserId    int
	ServiceId int
	Params    string
}

func ReceiveAndHandle() {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "kafka:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{"srvcs-topic1"}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			rqsrvcs := RequestedService{}
			err = json.Unmarshal([]byte(msg.Value), &rqsrvcs)
			rand.Seed(time.Now().UnixNano())
			result := rand.Float32() < 0.5
			if result == true {
				rqsrvcs.Params = "Выполнено успешно"
			} else {
				rqsrvcs.Params = "Не выполнено"
			}
			sendMessageToHis(rqsrvcs)
		} else {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
		}
	}

	c.Close()
}

func sendMessageToHis(rqSrvc RequestedService) {
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

	topic := "srvcs-topic2"
	for _, word := range []string{string(srvcString)} {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(word),
		}, nil)
	}
	p.Flush(15 * 1000)
}

func main(){
	ReceiveAndHandle()
}
