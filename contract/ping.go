package contract

import (
	"encoding/json"
	"time"
)

//PingModel represent team model
type PingModel struct {
	Url          string    `json:"url"`
	Start        time.Time `json:"start"`
	Connect      time.Time `json:"connect"`
	Dns          time.Time `json:"dns"`
	TlsHandShake time.Time `json:"handshake"`
	DnsDone      int64     `json:"dnsDone"`
	TlsDone      int64     `json:"tlsDone"`
	ConnectTime  int64     `json:"connectTime"`
	FirstByte    int64     `json:"firstByte"`
	//Headers      map[string][]string `json:"headers"`
}

type UrlPayload struct {
	URL string `json:"url"`
}

type QuePing struct {
	URL      string    `json:"url"`
	PingTime time.Time `json:"ping_time"`
}

func (q QuePing) Marshal() []byte {
	b, _ := json.Marshal(q)
	return b
}

func NewQuePing(payload []byte) (QuePing, error) {
	var q QuePing
	err := json.Unmarshal(payload, &q)
	return q, err
}
