package broker

import (
	"log"
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)


func (b *Broker) PublishMessage(sr StockReponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body, err := json.Marshal(sr)
	if err != nil {
		log.Printf("Response structure error %s ", err)
	}

	err = b.Channel.PublishWithContext(ctx,
		"",                    // exchange
		b.ProducerQueue.Name, // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})

	if err != nil {
		log.Fatalf("Failed to publish message: %s", err)
	}
}