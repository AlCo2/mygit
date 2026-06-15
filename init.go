package main

import (
	"fmt"
	"os"
)

var Folders = []string{
	MAIN_PATH,
	MAIN_PATH + "/objects",
	MAIN_PATH + "/objects/pack",
	MAIN_PATH + "/objects/info",
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func createFolders() error {
	for _, folder := range Folders {
		err := os.Mkdir(folder, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func Init() {

	if folderExists(".mygit") {
		fmt.Println("mygit already initialised")
		return
	}

	err := createFolders()

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
