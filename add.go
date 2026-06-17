package main

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"main/helper"
	"main/index"
	"os"
)

type Index struct {
	Header  index.Header
	Entries []index.Entry
}

// func parseIndex() Index {

// }

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

func addAllFiles(path string) {

}
func addToIndex(path string) error {
	files := os.Args[3:]

	if files[0] == "." {
		addAllFiles(path)
	}
	// todo later...
	return nil
}

func Add() {
	if len(os.Args) < 3 {
		fmt.Print("Nothing specified, nothing added.\nhint: Maybe you wanted to say 'mygit add .'?")
		return
	}

	path, err := helper.LocateMygitPath()

	if err != nil {
		fmt.Println("fatal: not a git repository (or any of the parent directories): .mygit")
		return
	}

	err = createIndex(path)

	if err != nil {
		fmt.Println("Error to create index file")
		return
	}

	err = addToIndex(path)
	if err != nil {
		fmt.Println("failed to add files")
		return
	}
}
