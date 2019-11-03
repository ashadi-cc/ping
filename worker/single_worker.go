package worker

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"ping/broker"
	"ping/contract"
	"time"

	"github.com/streadway/amqp"
)

var singleConfigBroker = broker.Config{
	ExchangeName: "single-worker-ping",
	Kind:         broker.DirectKind,
	Route:        "single-route-ping",
}

//SingleWorker hold worker instnce
type SingleWorker struct {
	messageBroker *broker.Broker
	WorkerId      int
}

//NewSingleWorker instance worker
func NewSingleWorker(workerId int) *SingleWorker {
	msgBroker := broker.NewBroker(configBroker)
	return &SingleWorker{
		messageBroker: msgBroker,
		WorkerId:      workerId,
	}
}

//Listen start listening message broker
func (w *SingleWorker) Listen() {
	fn := func(msgs <-chan amqp.Delivery) {
		for d := range msgs {
			que, err := contract.NewQuePing(d.Body)
			if err != nil {
				log.Printf("Error unmarshal message to Que Model: %s, %s", string(d.Body), err.Error())
			} else {
				w.Ping(que)
			}
		}
	}
	err := w.messageBroker.Consume(fn)
	if err != nil {
		log.Println("error when listening from message broker:", err.Error())
	}
}

//Ping ping url
func (w *SingleWorker) Ping(que contract.QuePing) {
	//if ping time less than now then republish
	if !time.Now().After(que.PingTime) {
		//sleep 100 milisecond before publish it
		time.Sleep(time.Millisecond * 100)
		w.Publish(que)
		return
	}
	w.ProcessUrl(que)
}

func (w *SingleWorker) ProcessUrl(que contract.QuePing) {
	log.Println("Worker:", w.WorkerId, "Processing URL:", que.URL)
	data := contract.UrlPayload{URL: que.URL}
	payload, _ := json.Marshal(data)
	response, err := http.Post("http://localhost:8080/ping", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error when call api ping service:", err.Error())
		return
	}
	//print response
	printResponse(response)
	//modify ping time with 10 second from now
	que.PingTime = time.Now().Add(time.Second * 10)
	//sleep 100 milisecond before publish it
	time.Sleep(time.Millisecond * 100)
	//publish url
	w.Publish(que)
}

//Publish resend back url to message broker
func (w *SingleWorker) Publish(que contract.QuePing) {
	b := que.Marshal()
	err := w.messageBroker.Publish(b)
	if err != nil {
		log.Println("Error to publish url to broker:", err.Error())
		return
	}
}
