package main

import (
	"fmt"
	"main/helper"
	"os"
	"path/filepath"
)

func createFolders(folders []string) error {
	for _, folder := range folders {
		if helper.FolderExists(folder) {
			continue
		}

		err := os.Mkdir(folder, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func setupFolders(path string) error {
	var folders = []string{
		path,
		path + "/objects",
		path + "/objects/pack",
		path + "/objects/info",
	}

	if len(os.Args) > 2 && !helper.FolderExists(os.Args[2]) {
		folders = append([]string{os.Args[2]}, folders...)
	}

	return createFolders(folders)
}

func Init() {
	path := ".mygit"
	if len(os.Args) > 2 {
		path = fmt.Sprintf("%s/.mygit", os.Args[2])
	}

	if helper.FolderExists(path) {
		fmt.Println("mygit already initialised")
		return
	}

	err := setupFolders(path)
	if err != nil {
		fmt.Println("ERROR: Failed to init mygit")
		return
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("ERROR: failed to read current directory")
	}

	fmt.Printf("Initialized empty Mygit repository in %s", absPath)
}
