package rabbitmq

import "testing"

const host = "amqp://test:test@192.168.1.5:5672"

func TestRabbitMQ_Publish(t *testing.T) {
	queue := New(host)
	defer queue.Close()
	queue.Bind("test")
}
