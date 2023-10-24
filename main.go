package main

import (
	"log"

	"github.com/starlingilcruz/golang-broker/broker"
)

var br broker.Broker

func main() {
	log.Println("Starting service ...")

	conn, ch := br.Connect()
	defer conn.Close()
	defer ch.Close()

	br.SetUp()

	go br.ReadMessages()
	select {}
}