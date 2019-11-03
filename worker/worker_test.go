package worker

import (
	"ping/contract"
	"testing"
	"time"
)

func TestMainWorker(t *testing.T) {
	list := []string{
		"https://www.google.com",
		"https://amazon.com",
		"https://facebook.com",
		"https://kompas.com",
		"https://detik.com",
	}
	worker := NewWorker(1)
	for _, url := range list {
		que := contract.QuePing{
			URL:      url,
			PingTime: time.Now(),
		}
		worker.Publish(que)
	}

}
