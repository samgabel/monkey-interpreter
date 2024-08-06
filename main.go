package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/samgabel/monkey-interpreter/repl"
)

func main() {
	// grab the OS user
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	// handle the startup message
	fmt.Printf("Hello %s! This is the Monkey programming language!\n", user.Username)
	fmt.Printf("Feel free to type in commands\n")

	// start our REPL
	repl.Start(os.Stdin, os.Stdout)
}
