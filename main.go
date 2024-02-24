// main.go
package main

import (
	"fmt"
	"main/server"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Iteringo programming language!\n",
		user.Username)
	fmt.Printf("Type `help()` in commands for common guide\n")
	// repl.Start(os.Stdin, os.Stdout)
	server.Start()
}
