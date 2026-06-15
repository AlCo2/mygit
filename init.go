package main

import (
	"encoding/binary"
	"fmt"
	"main/index"
	"os"
	"path/filepath"
)

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func createFolders(folders []string) error {
	for _, folder := range folders {
		if folderExists(folder) {
			continue
		}

		err := os.Mkdir(folder, 0755)

		if err != nil {
			return err
		}
	}

	return nil
}

func createIndex(path string) error {
	index_path := path + "/index"

	header := index.Header{
		Signature:  [4]byte{'D', 'I', 'R', 'C'},
		Version:    2,
		EntryCount: 0,
	}

	if fileExists(index_path) {
		return nil
	}

	file, err := os.Create(index_path)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := binary.Write(file, binary.BigEndian, header); err != nil {
		return err
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

	if len(os.Args) > 2 && !folderExists(os.Args[2]) {
		folders = append([]string{os.Args[2]}, folders...)
	}

	return createFolders(folders)
}

func Init() {
	path := ".mygit"
	if len(os.Args) > 2 {
		path = fmt.Sprintf("%s/.mygit", os.Args[2])
	}

	if folderExists(path) {
		fmt.Println("mygit already initialised")
		return
	}

	err := setupFolders(path)
	if err != nil {
		fmt.Println("ERROR: Failed to init mygit")
		return
	}

	err = createIndex(path)

	if err != nil {
		fmt.Println("ERROR: Failed to Create Index")
		return
	}

	absPath, err := filepath.Abs(path)

	if err != nil {
		fmt.Println("ERROR: failed to read current directory")
	}

	fmt.Printf("Initialized empty Mygit repository in %s", absPath)
}
