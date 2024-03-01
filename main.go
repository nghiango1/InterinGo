package main

import (
	"flag"
	"fmt"
	"interingo/repl"
	"interingo/server"
	"interingo/share"
	"os"
	"os/user"
)

// Run mode - Server flag
var verboseMode bool

// Run mode - Server flag
var serverMode bool

// Address we listen on.
var listenAddress string

// Address we listen on.
var fileLocation string

func init() {
	const (
		defaultServerMode    = false
		defaultVerboseMode    = false
		defaultListenAddress = "0.0.0.0:8080"
		defaultFileLocation  = ""
		serverUsage          = "Start as server mode"
		verboseUsage          = "Start as verbose mode, InterinGo will print a lot more infomation for Lexer, Parse and Evaluation product"
		listenAdrUsage       = "Listen address"
		fileLocationUsage       = "Using a file as input to parse, default as \"\" which mean not using file input"
	)

	flag.BoolVar(&serverMode, "server", defaultServerMode, serverUsage)
	flag.BoolVar(&serverMode, "s", defaultServerMode, serverUsage+" (shorthand)")
	flag.StringVar(&listenAddress, "listen-address", defaultListenAddress, listenAdrUsage)
	flag.StringVar(&listenAddress, "l", defaultListenAddress, listenAdrUsage+" (shorthand)")
	flag.StringVar(&fileLocation, "file", defaultFileLocation, fileLocationUsage)
	flag.StringVar(&fileLocation, "f", defaultFileLocation, fileLocationUsage+" (shorthand)")
	flag.BoolVar(&verboseMode, "verbose", defaultVerboseMode, verboseUsage)
	flag.BoolVar(&verboseMode, "v", defaultVerboseMode, verboseUsage+" (shorthand)")
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if verboseMode {
		fmt.Println("Verbose mode enable")
		share.VerboseMode = verboseMode
	}

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
