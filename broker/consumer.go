package broker

import (
	"log"
)

func (b *Broker)ReadMessages() {
	msgs, err := b.Channel.Consume(
		b.ConsumerQueue.Name, // queue
		"",                   // consumer
		true,                 // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	rsvdMsgs := make(chan StockRequest)
	go messageTransformer(msgs, rsvdMsgs)
	go processRequest(rsvdMsgs, b)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
}