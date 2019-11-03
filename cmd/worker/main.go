package main

import (
	"log"
	"ping/worker"
)

func main() {
	log.Println("Running worker")
	worker.RunService()
}
