package api

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"ping/contract"
	"sync"
	"time"

	"github.com/pkg/errors"
)

const maxPing = 10

var currentPing = 0

var m sync.Mutex

func PingUrl(url string) (*contract.PingModel, error) {
	m.Lock()
	if currentPing > maxPing {
		m.Unlock()
		return nil, fmt.Errorf("max allowed ping request reached. rejected")
	}

	currentPing = currentPing + 1
	m.Unlock()

	//subtract current ping for any return condition
	defer func() {
		m.Lock()
		currentPing = currentPing - 1
		m.Unlock()
	}()

	//spesify request timeouts by default
	http.DefaultClient.Timeout = time.Minute * 3
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "problem when ping url %s", url)
	}

	var start, connect, dns, tlsHandshake time.Time
	var dnsDone, tlsDone, connectTime, firstByte int64

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			dnsDone = time.Since(dns).Milliseconds()
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tlsDone = time.Since(tlsHandshake).Milliseconds()
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			connectTime = time.Since(connect).Milliseconds()
		},

		GotFirstResponseByte: func() {
			firstByte = time.Since(start).Milliseconds()
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()

	_, err = http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, errors.Wrapf(err, "problem when get response from url %s", url)
	}

	pingModel := &contract.PingModel{
		Url:          url,
		Start:        start,
		Connect:      connect,
		Dns:          dns,
		TlsHandShake: tlsHandshake,
		DnsDone:      dnsDone,
		TlsDone:      tlsDone,
		ConnectTime:  connectTime,
		FirstByte:    firstByte,
		//Headers:      res.Header,
	}

	return pingModel, err
}
