package main

import (
	"flag"
	"fmt"
	"interingo/pkg/repl"
	"interingo/pkg/server"
	"interingo/pkg/share"
	"log/slog"
	"os"
	"os/user"
)

// Runtime
var verboseMode bool

// Server mode flag
var serverMode bool

// Address we listen on - Use in Server mode flag
var listenAddress string

// Input file directory - Use in File mode flag
var fileLocation string

// Default level - Can be overide by make file
var debugLever slog.Level

func init() {
	const (
		defaultServerMode    = false
		defaultVerboseMode   = false
		defaultListenAddress = "0.0.0.0:8080"
		defaultFileLocation  = ""
		serverUsage          = "Start as server mode"
		verboseUsage         = "Start as verbose mode, InterinGo will print a lot more dev debug log"
		listenAdrUsage       = "Listen address"
		fileLocationUsage    = "Using a file as input to parse, default as \"\" which mean not using file input"
	)

	flag.BoolVar(&serverMode, "server", defaultServerMode, serverUsage)
	flag.BoolVar(&serverMode, "s", defaultServerMode, serverUsage+" (shorthand)")
	flag.StringVar(&listenAddress, "listen-address", defaultListenAddress, listenAdrUsage)
	flag.StringVar(&listenAddress, "l", defaultListenAddress, listenAdrUsage+" (shorthand)")
	flag.StringVar(&fileLocation, "file", defaultFileLocation, fileLocationUsage)
	flag.StringVar(&fileLocation, "f", defaultFileLocation, fileLocationUsage+" (shorthand)")
	flag.BoolVar(&verboseMode, "verbose", defaultVerboseMode, verboseUsage)
	flag.BoolVar(&verboseMode, "v", defaultVerboseMode, verboseUsage+" (shorthand)")

	flag.Parse()
	if verboseMode {
		debugLever = slog.LevelDebug
	}
}

func main() {
	share.SetDefaultLog(debugLever)
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if verboseMode {
		slog.Debug("Verbose mode enable", "level", debugLever)
	}

	if fileLocation != "" {
		fileContent, err := os.ReadFile(fileLocation)
		if err != nil {
			slog.Error("File read error, recheck file location", "error", err)
			return
		}
		repl := repl.NewRepl(nil, os.Stdin, os.Stdout)
		repl.Handle(string(fileContent))
		return
	}

	fmt.Printf("Hello %s! This is the Iteringo programming language!\n",
		user.Username)

	if serverMode {
		s := server.NewServer()
		s.Start(listenAddress)
		return
	}

	fmt.Printf("Type `help()` in commands for common guide\n")
	repl := repl.NewRepl(nil, os.Stdin, os.Stdout)
	repl.Start()
}
