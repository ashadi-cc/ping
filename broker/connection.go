package broker

import (
	"sync"

	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var initConn sync.Once

func Getconnection() (*amqp.Connection, error) {
	var err error
	initConn.Do(func() {
		conn, err = amqp.Dial("amqp://user:user@localhost:5672/")
	})

	return conn, err
}
