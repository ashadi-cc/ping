package broker

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/streadway/amqp"
)

func TestConnection(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	_, err := Getconnection()
	if err != nil {
		t.Fatalf("cannot create connection, got error %s", err.Error())
	}
}

var config = Config{
	ExchangeName: "test-exchange-name",
	Kind:         DirectKind,
	Route:        "test-route",
}

var broker *Broker

func TestPublisher(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	broker = NewBroker(config)
	err := broker.Publish([]byte("test-data"))
	if err != nil {
		t.Fatalf("cannot publish message %s", err.Error())
	}
}

func TestConsume(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	fn := func(msgs <-chan amqp.Delivery) {
		wt, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		go func() {
			for d := range msgs {
				log.Printf("Received message : %s \n", string(d.Body))
			}
		}()
		//hold until timeout
		<-wt.Done()
	}

	err := broker.Consume(fn)
	if err != nil {
		t.Fatalf("cannot consume: %s", err.Error())
	}
}
