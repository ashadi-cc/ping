package worker

import "log"

//RunService get worker service
func RunService() {
	forever := make(chan bool)
	workers := 10
	log.Printf("Create %d Workers \n", workers)
	createWorker(workers)
	<-forever
}

func createWorker(workers int) {
	for i := 1; i <= workers; i++ {
		log.Printf("Worker :%d created \n", i)
		go func(id int) {
			//worker with only single process
			//NewSingleWorker(id).Listen()

			//worker with multi process
			NewWorker(id).Listen()
		}(i)
	}
}
