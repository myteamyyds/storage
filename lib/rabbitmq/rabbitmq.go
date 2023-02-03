package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type RabbitMQ struct {
	chanel   *amqp.Channel
	conn     *amqp.Connection
	Name     string
	exchange string
}

func New(s string) *RabbitMQ {
	connection, err := amqp.Dial(s)
	if err != nil {
		log.Fatal(err)
	}
	channel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}
	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}
	mq := new(RabbitMQ)
	mq.chanel = channel
	mq.conn = connection
	mq.Name = queue.Name
	return mq
}

func (q *RabbitMQ) Bind(exchange string) {
	err := q.chanel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	q.exchange = exchange
}

func (q *RabbitMQ) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	err = q.chanel.Publish("", queue, false, false, amqp.Publishing{
		ReplyTo: q.Name,
		Body:    []byte(str),
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (q *RabbitMQ) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	err = q.chanel.Publish(
		exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body:    []byte(str),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}

func (q *RabbitMQ) Consume() <-chan amqp.Delivery {
	consume, err := q.chanel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	return consume
}

func (q *RabbitMQ) Close() {
	q.chanel.Close()
	q.conn.Close()
}
