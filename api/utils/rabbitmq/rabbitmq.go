package rabbitmq

import (
	"os"

	"github.com/streadway/amqp"
)

func ConnectMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://" + os.Getenv("RABBITMQ_DEFAULT_USER") + ":" + os.Getenv("RABBITMQ_DEFAULT_PASS") + "@" + os.Getenv("RABBITMQ_DEFAULT_HOST") + ":5672/")
	return conn, err
}
