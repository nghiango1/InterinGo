// main.go
package main

import (
	"flag"
	"fmt"
	"main/repl"
	"main/server"
	"os"
	"os/user"
)

// Run mode - Server flag
var serverMode bool

// Address we listen on.
var listenAddress string = ":8080"

func init() {
	const (
		defaultServerMode    = false
		defaultListenAddress = "127.0.0.1:8080"
		serverUsage          = "Start as server mode"
		listenAdrUsage       = "Listen address"
	)

	flag.BoolVar(&serverMode, "server", defaultServerMode, serverUsage)
	flag.BoolVar(&serverMode, "s", defaultServerMode, serverUsage+" (shorthand)")
	flag.StringVar(&listenAddress, "listen-address", defaultListenAddress, listenAdrUsage)
	flag.StringVar(&listenAddress, "l", defaultListenAddress, listenAdrUsage+" (shorthand)")
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Iteringo programming language!\n",
		user.Username)

	flag.Parse()
	if serverMode {
		server.Start(listenAddress)
	} else {
		fmt.Printf("Type `help()` in commands for common guide\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}
