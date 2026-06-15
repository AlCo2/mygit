package main

import (
	"fmt"
	"os"
)

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func Init() {

	if folderExists(".mygit") {
		fmt.Println("mygit already initialised")
		return
	}

	err := os.Mkdir(".mygit", 0755)

	if err != nil {
		fmt.Println("ERROR: Failed to init mygit")
		return
	}
	dir, err := os.Getwd()

	if err != nil {
		fmt.Println("ERROR: failed to read current directory")
	}

	fmt.Printf("Initialized empty Git repository in %s", dir)
}
