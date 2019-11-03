package broker

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

//Kind type definition
type Kind string

const (
	DirectKind Kind = "direct"
	TopicKind  Kind = "topic"
	FanoutKind Kind = "fanout"
)

//BrokerConfig hold config variable for broker
type Config struct {
	ExchangeName string
	Kind         Kind
	Route        string
}

//Broker represent broker instance
type Broker struct {
	Conn   *amqp.Connection
	Config Config
}

//NewBroker new instance of broker
func NewBroker(config Config) *Broker {
	broker := &Broker{Config: config}
	broker.Init()
	return broker
}

//Init create connection
func (b *Broker) Init() {
	conn, err := Getconnection()
	if err != nil {
		log.Fatalf("Can not initalize to rabbitmq server: %s", err.Error())
	}
	b.Conn = conn
}

//Publish publish message to rabbitmq
func (b *Broker) Publish(data []byte) error {
	channel, err := b.Conn.Channel()
	if err != nil {
		return errors.Wrap(err, "Can not create channel")
	}
	defer channel.Close()

	_, err = b.declareAndBind(channel)
	if err != nil {
		return err
	}

	message := amqp.Publishing{
		Body: data,
	}
	err = channel.Publish(b.Config.ExchangeName, b.Config.Route, false, false, message)
	if err != nil {
		return errors.Wrap(err, "cannot publish message")
	}

	return nil
}

//ConsumeFunc hold func
type ConsumeFunc func(msgs <-chan amqp.Delivery)

//Consume consume message
func (b *Broker) Consume(consumeFn ConsumeFunc) error {
	channel, err := b.Conn.Channel()
	if err != nil {
		return errors.Wrap(err, "Can not create channel")
	}
	defer channel.Close()

	queName, err := b.declareAndBind(channel)
	if err != nil {
		return err
	}

	msgs, err := channel.Consume(queName, "", true, false, false, false, nil)
	if err != nil {
		return errors.Wrap(err, "Cannot consume channel")
	}

	//consume
	consumeFn(msgs)

	return nil
}

func (b *Broker) declareAndBind(channel *amqp.Channel) (string, error) {
	queName := fmt.Sprintf("que-%s-%s", b.Config.Kind, b.Config.Route)

	// read this for type of kind object https://medium.com/faun/different-types-of-rabbitmq-exchanges-9fefd740505d
	err := channel.ExchangeDeclare(b.Config.ExchangeName, string(b.Config.Kind), true, false, false, false, nil)
	if err != nil {
		return queName, errors.Wrapf(err, "Cannot declare exchange %+v", b.Config)
	}

	que, err := channel.QueueDeclare(queName, true, false, false, false, nil)
	if err != nil {
		return queName, errors.Wrapf(err, "cannot declare queue: %s", queName)
	}

	err = channel.QueueBind(que.Name, b.Config.Route, b.Config.ExchangeName, false, nil)
	if err != nil {
		return queName, errors.Wrapf(err, "Cannot binding queue %+v", b.Config)
	}
	return queName, nil
}
