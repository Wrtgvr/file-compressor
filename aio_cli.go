package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// TODO: code refactoring
func aio_cli_start() {
	actionChosen := false
	for !actionChosen {
		var action int

		fmt.Print("Choose action:\n1. Compress files to .wrtaio\n2. Decompress .wrtaio file\n3. Quit\nType action number: ")
		_, scanErr := fmt.Scan(&action) // scanErr handles in switch-case default.
		// set to true here so theres no need to set it to true in every action case.
		// it will be set back to false in switch-case default.
		actionChosen = true

		switch action {
		case 1:
			action_compress()
		case 2:
			action_decompress()
		case 3:
			os.Exit(0)
		default:
			actionChosen = false
			if scanErr != nil {
				log.Fatalf("ERROR: %v", scanErr)
			}
			fmt.Println("Unknown action")
		}
	}
}

func action_decompress() {
	//* get files in input folder
	inputFiles, err := os.ReadDir(INPUT_FOLDER)
	if err != nil {
		log.Fatalf("Unable to get files in input folder: %v", err)
	}

	//* get files with extension == FILE_EXTENSION
	compressedFiles := map[int]string{}
	count_compressedFiles := 1

	for _, ent := range inputFiles {
		s := strings.Split(ent.Name(), ".")
		if s[len(s)-1] == FILE_EXTENSION {
			compressedFiles[count_compressedFiles] = ent.Name()
			count_compressedFiles++
		}
	}

	//* check if theres no compressed files
	if len(compressedFiles) == 0 {
		var a string
		fmt.Printf("Theres no compressed files in input folder.\nCompressed files has extension .%s", FILE_EXTENSION)
		fmt.Scan(&a)
		os.Exit(0)
	}

	//* print list of files which can be decompressed to user
	Message_FilesToDeompress := "Files to decompress:\n"

	for i, filename := range compressedFiles {
		Message_FilesToDeompress += fmt.Sprintf("%d. %s\n", i, filename)
	}
	fmt.Println(Message_FilesToDeompress)

	//* let used choose file to decompress
	var isFileChosen bool
	var chosenFile string

	for !isFileChosen {
		var fileNum int
		fmt.Print("Choose file to decompress\nType file number: ")
		fmt.Scan(&fileNum)

		f, ok := compressedFiles[fileNum]
		if !ok {
			fmt.Println("Invalid file number.\nIf you trying to quit then choose any file and decline confirmation during next step.")
			continue
		}

		isFileChosen = true
		chosenFile = f
	}

	fmt.Printf("Chosen file: %s\n", chosenFile)

	//* get confirmation
	var actionChosen bool

	for !actionChosen {
		var confirmation string

		fmt.Print("Decompress chosen file?\n(Y/N): ")
		fmt.Scan(&confirmation)

		actionChosen = true

		switch strings.ToLower(confirmation) {
		case "y":
			decompressWRTAIO(fmt.Sprintf("%s/%s", INPUT_FOLDER, chosenFile))
		case "n":
			os.Exit(0)
		default:
			actionChosen = false
		}
	}
}

func action_compress() {
	//* get files in input folder
	inputFiles, err := os.ReadDir(INPUT_FOLDER)
	if err != nil {
		log.Fatalf("Unable to get files in input folder: %v", err)
	}

	//* check if theres no files to compress
	if len(inputFiles) == 0 {
		var a string
		fmt.Print("Theres no files in input folder")
		fmt.Scan(&a)
		os.Exit(0)
	}

	//* print list of files to compress to user
	Count_FilesToCompress := 1
	Message_FilesToCompress := "Files to compress:\n"

	for _, ent := range inputFiles {
		Message_FilesToCompress += fmt.Sprintf("%d. %s\n", Count_FilesToCompress, ent.Name())
		Count_FilesToCompress++
	}
	fmt.Println(Message_FilesToCompress)

	//* get confirmation
	var actionChosen bool

	for !actionChosen {
		var confirmation string

		fmt.Println("Compress listed files?\n(Y/N): ")
		fmt.Scan(&confirmation)

		actionChosen = true

		switch strings.ToLower(confirmation) {
		case "y":
			compressWRTAIO(INPUT_FOLDER)
		case "n":
			os.Exit(0)
		default:
			actionChosen = false
		}
	}
}
