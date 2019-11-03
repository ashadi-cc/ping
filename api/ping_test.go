package api

import (
	"encoding/json"
	"testing"
)

func TestPing(t *testing.T) {
	ping, err := PingUrl("https://www.google.com")
	if err != nil {
		t.Fatal(err)
	}
	_, err = json.Marshal(ping)
	if err != nil {
		t.Fatal("error when marshall to json", err)
	}

}
