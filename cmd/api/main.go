package main

import (
	"log"
	"ping/api"
)

func main() {
	log.Println("starting ping api service")
	api.InitServer()
}
