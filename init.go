package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"main/index"
	"os"
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
		Version:    1,
		EntryCount: 0,
	}

	index_file, err := os.Create(index_path)
	if err != nil {
		return err
	}

	defer index_file.Close()

	if err := binary.Write(index_file, binary.BigEndian, header); err != nil {
		return err
	}

	files, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				return err
			}

			hash := sha1.New()

			file_open, err := os.Open(file.Name())
			if err != nil {
				return err
			}

			if _, err := io.Copy(hash, file_open); err != nil {
				file_open.Close()
				return err
			}

			file_open.Close()
			entry := index.Entry{
				Size: uint32(info.Size()),
				SHA1: [20]byte(hash.Sum(nil)),
			}

			copy(entry.Path[:], info.Name())

			if err := binary.Write(index_file, binary.BigEndian, entry); err != nil {
				return err
			}

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

	// if folderExists(path) {
	// 	fmt.Println("mygit already initialised")
	// 	return
	// }

	// err := setupFolders(path)
	// if err != nil {
	// 	fmt.Println("ERROR: Failed to init mygit")
	// 	return
	// }

	err := createIndex(path)

	if err != nil {
		fmt.Println("ERROR: Failed to Create Index")
		log.Fatal(err)
		return
	}

	// absPath, err := filepath.Abs(path)

	// if err != nil {
	// 	fmt.Println("ERROR: failed to read current directory")
	// }

	// fmt.Printf("Initialized empty Mygit repository in %s", absPath)
}
