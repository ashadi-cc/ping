package worker

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"ping/broker"
	"ping/contract"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const MaxParalelProcess = 10

var configBroker = broker.Config{
	ExchangeName: "worker-ping",
	Kind:         broker.DirectKind,
	Route:        "route-ping",
}

//Worker hold worker instnce
type Worker struct {
	mutex         sync.Mutex
	Urls          []contract.QuePing
	maxParalel    int
	workers       int
	messageBroker *broker.Broker
	WorkerId      int
}

//NewWorker instance worker
func NewWorker(workerId int) *Worker {
	msgBroker := broker.NewBroker(configBroker)
	return &Worker{
		messageBroker: msgBroker,
		maxParalel:    MaxParalelProcess,
		WorkerId:      workerId,
	}
}

//Listen start listening message broker
func (w *Worker) Listen() {
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

func (w *Worker) Ping(que contract.QuePing) {
	//if ping time less than now then republish
	if !time.Now().After(que.PingTime) {
		//sleep 100 milisecond before publish it
		<-time.Tick(time.Millisecond * 100)
		w.Publish(que)
		return
	}

	w.mutex.Lock()
	w.Urls = append(w.Urls, que)
	defer w.mutex.Unlock()
	//if length of url greater than max process then just wait it
	if len(w.Urls) > w.maxParalel {
		return
	}
	w.workers = w.workers + 1
	//run concurrent process
	go w.Run(w.workers)
}

func (w *Worker) Run(workerId int) {
	w.mutex.Lock()
	length := len(w.Urls)

	//if length is zero just return
	if length == 0 {
		w.mutex.Unlock()
		return
	}

	//get first url
	url := w.Urls[0]
	//process url
	w.Urls = w.Urls[1:]
	//unlock
	w.mutex.Unlock()

	//process in background with delay
	w.ProcessUrl(url, workerId)

	//process remain url in list
	w.Run(workerId)
}

func (w *Worker) ProcessUrl(que contract.QuePing, workerId int) {
	if w.workers > w.maxParalel {
		<-time.Tick(time.Millisecond * 100)
	}

	log.Println("Worker:", w.WorkerId, "SUB WORKER:", workerId, "Processing URL:", que.URL)
	data := contract.UrlPayload{URL: que.URL}
	payload, _ := json.Marshal(data)
	response, err := http.Post("http://localhost:8081/ping", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Println("Error when call api ping service:", err.Error())
		return
	}
	w.workers = w.workers - 1
	//print response
	printResponse(response)
	//modify ping time with 10 second from now
	que.PingTime = time.Now().Add(time.Second * 10)
	//sleep 100 milisecond before publish it
	<-time.Tick(time.Millisecond * 100)
	//publish url
	w.Publish(que)
}

//Publish resend back url to message broker
func (w *Worker) Publish(que contract.QuePing) {
	b := que.Marshal()
	err := w.messageBroker.Publish(b)
	if err != nil {
		log.Println("Error to publish url to broker:", err.Error())
		return
	}
}

func printResponse(r *http.Response) {
	data, _ := ioutil.ReadAll(r.Body)
	log.Println("RESULT:", string(data))
}
