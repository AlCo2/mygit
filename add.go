package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"main/helper"
	"main/index"
	"os"
)

type Index struct {
	Header    index.Header
	Entries   []index.Entry
	HashTable map[string]*index.Entry
	PathTable map[string]*index.Entry
}

func parseIndex() *Index {
	mygit_path, err := helper.LocateMygitPath()
	if err != nil {
		return nil
	}

	var header index.Header
	var entries []index.Entry

	index_file, err := os.ReadFile(mygit_path + "/index")

	if err != nil {
		return nil
	}
	r := bytes.NewReader(index_file)

	err = binary.Read(r, binary.BigEndian, &header)
	if err != nil {
		return nil
	}

	for i := uint32(0); i < header.EntryCount; i++ {
		var e index.Entry

		if err := binary.Read(r, binary.BigEndian, &e.Size); err != nil {
			return nil
		}

		if _, err := io.ReadFull(r, e.SHA1[:]); err != nil {
			return nil
		}

		if _, err := io.ReadFull(r, e.Path[:]); err != nil {
			return nil
		}

		entries = append(entries, e)
	}

	hashtTable := make(map[string]*index.Entry, header.EntryCount)
	pathTable := make(map[string]*index.Entry, header.EntryCount)

	for i := range entries {
		e := &entries[i]
		key := hex.EncodeToString(e.SHA1[:])
		path := string(e.Path[:])
		pathTable[path] = e
		hashtTable[key] = e
	}

	return &Index{
		Header:    header,
		Entries:   entries,
		HashTable: hashtTable,
		PathTable: pathTable,
	}
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

	return nil
}

func writeIndex(index Index) error {
	git_path, err := helper.LocateMygitPath()
	if err != nil {
		return err
	}

	index_path := git_path + "/index_"

	index_file, err := os.Create(index_path)
	if err != nil {
		index_file.Close()
		return err
	}

	if err := binary.Write(index_file, binary.BigEndian, index.Header); err != nil {
		index_file.Close()
		return err
	}

	for _, entry := range index.Entries {
		if err := binary.Write(index_file, binary.BigEndian, entry); err != nil {
			index_file.Close()
			return err
		}
	}

	err = os.Remove(git_path + "/index")
	if err != nil {
		index_file.Close()
		return err
	}

	index_file.Close()
	os.Rename(index_path, git_path+"/index")

	return nil
}

func addAllFiles(path string) error {
	index_path := path + "/index"

	index_file, err := os.Open(index_path)

	if err != nil {
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

func addFiles(files []string) error {
	current_index := parseIndex()

	if current_index == nil {
		return errors.New("error while parsing index file")
	}

	for _, file_path := range files {
		file, err := os.Stat(file_path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("fatal: pathspec '%s' did not match any files", file_path)
			}
			return err
		}

		old_entry := current_index.PathTable[file_path]

		if old_entry != nil {
			fmt.Println("todo: handle already added entry")
			return nil
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
			Size: uint32(file.Size()),
			SHA1: [20]byte(hash.Sum(nil)),
		}

		copy(entry.Path[:], file.Name())

		current_index.Header.EntryCount += 1
		current_index.Entries = append(current_index.Entries, entry)
	}

	writeIndex(*current_index)

	return nil
}

func addToIndex(path string) error {
	files := os.Args[2:]

	if files[0] == "." {
		err := addAllFiles(path)
		if err != nil {
			return err
		}
	}

	return addFiles(files)
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
		fmt.Println("failed to create index file")
		return
	}

	err = addToIndex(path)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
