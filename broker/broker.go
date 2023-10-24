package broker

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/joho/godotenv"

	"github.com/starlingilcruz/golang-broker/api"

	amqp "github.com/rabbitmq/amqp091-go"
)

type StockRequest struct {
	RoomName string `json:"roomName"`
	RoomId   uint   `json:"roomId"`
	Message  string `json:"message"`
}

type StockReponse struct {
	RoomId  uint   `json:"RoomId"`
	Message string `json:"Message"`
}

type Broker struct {
	Conn           *amqp.Connection
	ConsumerQueue  amqp.Queue
	ProducerQueue  amqp.Queue
	Channel        *amqp.Channel
}

func (b *Broker)Connect() (*amqp.Connection, *amqp.Channel) {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Failed to load .env: %s", err)
	}

	amqHost := os.Getenv("RMQ_HOST")
	amqUser := os.Getenv("RMQ_USERNAME")
	amqPass := os.Getenv("RMQ_PASSWORD")
	amqPort := os.Getenv("RMQ_PORT")

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", amqUser, amqPass, amqHost, amqPort))

	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Println("Failed to start a channel")
	}

	b.Conn    = conn
	b.Channel = ch

	return conn, ch
}

func (b *Broker)CreateQueue(queue string) (amqp.Queue, error) {
	q, err := b.Channel.QueueDeclare(
		queue,         // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)

	if err != nil {
		log.Println(fmt.Sprintf("Error: failed to create queue - %s", err))
	}

	return q, err
}

func (b *Broker) SetUp() {
	b.ConsumerQueue, _ = b.CreateQueue(os.Getenv("BR_CONSUMER_QUEUE"))
	b.ProducerQueue, _ = b.CreateQueue(os.Getenv("BR_PRODUCER_QUEUE"))
}

func messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan StockRequest) {
	var sr StockRequest
	for d := range entries {
		err := json.Unmarshal([]byte(d.Body), &sr)
		if err != nil {
			log.Printf("Error on received request : %s ", err)
			continue
		}
		receivedMessages <- sr
	}
}

func processRequest(s <-chan StockRequest, b *Broker) {

	for r := range s {
		cM := r.Message
		cM = strings.Replace(cM, "/stock=", "", 1)
		sr := StockReponse{
			RoomId:  r.RoomId,
			Message: fmt.Sprintf("Processing: %s", cM),
		}
		go b.PublishMessage(sr)
		msg := api.EvalStock(cM)
		sr2 := StockReponse{
			RoomId:  r.RoomId,
			Message: msg,
		}
		go b.PublishMessage(sr2)
	}
}