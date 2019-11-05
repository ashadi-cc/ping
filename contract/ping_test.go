package contract

import (
	"fmt"
	"testing"
	"time"
)

func TestQuePing(t *testing.T) {
	que := QuePing{
		PingTime: time.Now(),
	}
	fmt.Println(string(que.Marshal()))

	que.PingTime = que.PingTime.Add(time.Second * 10)

	fmt.Println(string(que.Marshal()))
}
