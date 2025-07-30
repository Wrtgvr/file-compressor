package main

import (
	"log"
	"os"
)

// var (
// 	// byte-separator, placed between file data and file name
// 	FILE_NAME_SEPARATOR_BYTES []byte = []byte("NAME_SEP")
// )

// 1. if symbol use N bits & N < 8, then we can use N bits instead of a byte to store that symbol
// (how to work with bits in golang???)

const (
	// how much bytes take for every iteration through file data
	FILE_DATA_ON_READ_PART_SIZE int64 = 200
	// instead of two bytes for line break program use one byte
	LINE_BREAK_COMPRESSED_BYTE byte = 0
	// byte-separator, placed between file data and file name
	FILE_NAME_SEPARATOR_BYTE byte = 255
	// folders paths
	INPUT_FOLDER  = "input"
	OUTPUT_FOLDER = "output"
	TEMP_FOLDER   = "temp"
	// Compressed file extansion
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
		os.Remove(INPUT_FOLDER)
		os.Mkdir(INPUT_FOLDER, 0644)
	}
	if !outputFolderExists {
		os.Remove(OUTPUT_FOLDER)
		os.Mkdir(OUTPUT_FOLDER, 0644)
	}
}

func main() {
	aio_cli_start()
}
