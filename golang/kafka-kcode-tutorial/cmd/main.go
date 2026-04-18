package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-kafka-kcode-tutorial/internal/consumer"
	"github.com/golang-kafka-kcode-tutorial/internal/producer"
)

type Server struct {
	producer *producer.KafkaProducer
	consumer *consumer.KafkaConsumer
	msgChan  chan string
}

func NewServer() *Server {
	msgChan := make(chan string, 64)
	return &Server{
		producer: producer.NewKafkaProducer(""),
		consumer: consumer.NewKafkaConsumer(msgChan),
		msgChan:  msgChan,
	}
}

func (serv *Server) ProduceMessage() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		msg := fmt.Sprintf("Message produced at %v", t)
		serv.producer.ProduceMessages(msg)
	}
}

// func (serv *Server) handleMessages() {
// 	for msg := range serv.msgChan {
// 		fmt.Printf("Handling message %s\n", msg)
// 	}
// }

func worker(wg *sync.WaitGroup, msgChan chan string, resChan chan string) {
	defer wg.Done()
	for msg := range msgChan {
		time.Sleep(time.Millisecond * 50)
		fmt.Printf("Message Processed : %s\n", msg)
		resChan <- msg
	}
}

func main() {
	var wg *sync.WaitGroup
	var totalWorkers int = 2
	resultChan := make(chan string, 64)

	server := NewServer()

	for range totalWorkers {
		wg.Add(1)
		go worker(wg, server.msgChan, resultChan)
	}

	for i := range resultChan {
		fmt.Printf("Result : %s\n", i)
	}

	go server.ProduceMessage()
}
