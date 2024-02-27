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
var listenAddress string

// Address we listen on.
var fileLocation string

func init() {
	const (
		defaultServerMode    = false
		defaultListenAddress = "0.0.0.0:8080"
		defaultFileLocation  = ""
		serverUsage          = "Start as server mode"
		listenAdrUsage       = "Listen address"
		fileLocationUseage       = "Using a file as input to parse, default as \"\" which mean not using file input"
	)

	flag.BoolVar(&serverMode, "server", defaultServerMode, serverUsage)
	flag.BoolVar(&serverMode, "s", defaultServerMode, serverUsage+" (shorthand)")
	flag.StringVar(&listenAddress, "listen-address", defaultListenAddress, listenAdrUsage)
	flag.StringVar(&listenAddress, "l", defaultListenAddress, listenAdrUsage+" (shorthand)")
	flag.StringVar(&fileLocation, "file", defaultFileLocation, fileLocationUseage)
	flag.StringVar(&fileLocation, "f", defaultFileLocation, fileLocationUseage+" (shorthand)")
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if fileLocation != "" {
		fileContent, err := os.ReadFile(fileLocation)
		if err != nil {
			fmt.Println("File read error, recheck file location, error code:", err)
			return
		} 
		repl.Handle(string(fileContent), os.Stdout)
		return
	}

	fmt.Printf("Hello %s! This is the Iteringo programming language!\n",
		user.Username)

	if serverMode {
		server.Start(listenAddress)
		return
	}

	fmt.Printf("Type `help()` in commands for common guide\n")
	repl.Start(os.Stdin, os.Stdout)
	
}
