package main

import (
	"fmt"
	"os"
)

const MAIN_PATH = ".mygit"

func main() {

	if len(os.Args) <= 1 {
		fmt.Println("USAGE: mygit <command>")
		return
	}

	command := os.Args[1]

	switch command {
	case "help":
		fmt.Println("soon...")
	case "status":
		Status()
	case "init":
		Init()
	default:
		fmt.Printf("mygit: '%s' is not a git command. See 'mygit help'.", command)
	}

}
