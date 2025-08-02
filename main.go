package main

import (
	"log"
	"os"
)

// 1. if symbol use N bits & N < 8, then we can use N bits instead of a byte to store that symbol
// I forgot name of the guy who invented method which use binary tree to separate data by amount of used bits
// Also I didn't really understood that method, I need read about it again.

const (
	// how much bytes take for every iteration through file data
	FILE_DATA_ON_READ_PART_SIZE int64 = 200 // tbh idk what size will be best
	//! I don't think changing FILE_..._BYTES_LENGTH is good idea
	//! cuz if it changed then you might need to change some lines like this:
	//! binary.BigEndian.PutUint16(...)
	//!! just don't change it unless you ready to search for every piece of code which need to be changed
	// how much bytes use for storing length of file name
	FILE_NAME_BYTES_LENGTH int8 = 2 // use N bytes to store file name/path length
	// how much bytes use for storing length of file data
	FILE_DATA_BYTES_LENGTH int8 = 4 // use N bytes to store file data length
	// folders paths
	INPUT_FOLDER  = "input"
	OUTPUT_FOLDER = "output"
	// Compressed file extension
	FILE_EXTENSION = "wrtaio"
)

func init() {
	entries, err := os.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	inputFolderExists := false
	outputFolderExists := false

	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}

		switch ent.Name() {
		case INPUT_FOLDER:
			inputFolderExists = true
		case OUTPUT_FOLDER:
			outputFolderExists = true
		}
	}

	if !inputFolderExists {
		os.Remove(INPUT_FOLDER) // is theres file named "input" which is NOT a dir
		os.Mkdir(INPUT_FOLDER, 0644)
	}
	if !outputFolderExists {
		os.Remove(OUTPUT_FOLDER) // is theres file named "output" which is NOT a dir
		os.Mkdir(OUTPUT_FOLDER, 0644)
	}
}

func main() {
	aio_cli_start()
}
