package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"ping/contract"
	"ping/worker"
	"time"
)

func main() {
	file, err := os.Open("web.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	worker := worker.NewWorker(1)

	scanner := bufio.NewScanner(file)
	now := time.Now()
	for scanner.Scan() {
		now = now.Add(time.Millisecond * 50)
		url := fmt.Sprintf("http://%s", scanner.Text())
		que := contract.QuePing{
			URL:      url,
			PingTime: now,
		}
		log.Println(url)
		worker.Publish(que)
	}
}
