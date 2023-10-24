package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/starlingilcruz/golang-broker/broker"
	"github.com/starlingilcruz/golang-broker/utils"
	amqp "github.com/rabbitmq/amqp091-go"
)

var br broker.Broker

func main() {
	log.Println("Stock bot service starting ...")

	err := godotenv.Load()
	utils.FailOnError(err, "Error loading .env file")

	rmqHost := os.Getenv("RMQ_HOST")
	rmqUserName := os.Getenv("RMQ_USERNAME")
	rmqPassword := os.Getenv("RMQ_PASSWORD")
	rmqPort := os.Getenv("RMQ_PORT")
	dsn := "amqp://" + rmqUserName + ":" + rmqPassword + "@" + rmqHost + ":" + rmqPort + "/"

	log.Println(dsn)

	conn, err := amqp.Dial(dsn)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	br.SetUp(ch)
	go br.ReadMessages()
	select {}
}