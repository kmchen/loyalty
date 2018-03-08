package main

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) error {
	errMsg := fmt.Sprintf("%s: %s", msg, err)
	return errors.New(errMsg)
}

func NewQueueConn(addr string) (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	if conn, err = amqp.Dial(addr); err != nil {
		return nil, failOnError(err, "Failed to connect to RabbitMQ")
	}
	return conn, nil
}

func QueueChan(conn *amqp.Connection) (*amqp.Channel, error) {
	var ch *amqp.Channel
	var err error
	if ch, err = conn.Channel(); err != nil {
		return nil, failOnError(err, "Failed to open a channel")
	}

	err = ch.ExchangeDeclare(
		"events", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, failOnError(err, "Failed to declare an exchange")
	}
	return ch, nil
}

func Bind(ch *amqp.Channel, routingKeys []string, exchange string) (<-chan amqp.Delivery, error) {

	if len(routingKeys) == 0 {
		return nil, fmt.Errorf("Missing routing keys")
	}

	var q amqp.Queue
	var err error
	q, err = ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, failOnError(err, "Failed to declare a queue")
	}

	for _, val := range routingKeys {
		err = ch.QueueBind(
			q.Name, // queue name
			val,
			exchange, // exchange
			false,
			nil)
		if err != nil {
			return nil, failOnError(err, "Failed to bind a queue")
		}
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return nil, failOnError(err, "Failed to register a consumer")
	}
	return msgs, nil
}
