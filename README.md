# Distributed Stock Service

## Requirements

- Go
- RabbitMQ

### Run development: 

Make sure that you have the RabbitMQ service running and configure the environment variables before running the application.

Here is a sample of the environment variables:

RMQ_USERNAME = rabbitmq
RMQ_PASSWORD = rabbitmq
RMQ_HOST = localhost
RMQ_PORT = 5672
BR_CONSUMER_QUEUE = stockbot-receiver
BR_PRODUCER_QUEUE = stockbot-publisher

## Start Service

```
  go run main.go
```