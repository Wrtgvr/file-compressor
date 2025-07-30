package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// TODO: code refactoring
func decompressWRTAIO(filepath string) int64 {
	//* open file
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("decompress: unable to open file. does file exists? error: %w", err))
	}
	defer file.Close()

	decompressionStartTime := time.Now()

	//? map fileName = fileData
	decompData := map[string][]byte{}

	var EOFError bool
	var offset int64
	var isFileName bool
	var currentFileName string

	for !EOFError {
		readData := make([]byte, FILE_DATA_ON_READ_PART_SIZE)

		bytesRead, err := file.ReadAt(readData, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
			} else {
				log.Fatal(fmt.Errorf("decompress: unable to read file bytes at %d. error: %w", offset, err))
			}
		}
		offset += FILE_DATA_ON_READ_PART_SIZE

		for i := 0; i < bytesRead; i++ {
			if readData[i] == FILE_NAME_SEPARATOR_BYTE {
				isFileName = !isFileName
				if isFileName {
					currentFileName = ""
				}
				continue
			}
			if isFileName {
				currentFileName += string(readData[i])
				continue
			}

			// check if is new line
			if readData[i] == LINE_BREAK_COMPRESSED_BYTE {
				decompData[currentFileName] = append(decompData[currentFileName], 13, 10)
				continue
			}

			decompData[currentFileName] = append(decompData[currentFileName], readData[i])
		}
	}

	for n := range decompData {
		fmt.Println(n)
	}

	for filename, bytes := range decompData {
		file, err := os.Create(fmt.Sprintf("%s/%s", OUTPUT_FOLDER, filename))
		log.Println(filename)
		if err != nil {
			log.Fatal(fmt.Errorf("decompress: unable to create files with decompressed data. error: %w", err))
		}
		file.Write(bytes)
	}

	return time.Since(decompressionStartTime).Milliseconds()
}

func compressWRTAIO(folderName string) ([]byte, int64) {
	//* get folder files
	inputFiles, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to read directory. error: %w", err))
	}

	compressionStartTime := time.Now()

	folderData := []byte{}

	for _, dirEntry := range inputFiles {
		var data []byte
		if dirEntry.IsDir() {
			// if file is a folder then compress files from inner folder
			data, _ = compressWRTAIO(fmt.Sprintf("%s/%s", folderName, dirEntry.Name()))
		} else {
			//* file open
			file, err := os.Open(fmt.Sprintf("%s/%s", folderName, dirEntry.Name()))
			if err != nil {
				log.Print(fmt.Errorf("compress: unable to open file: %s. error: %w", dirEntry.Name(), err))
				continue
			}

			var EOFError bool
			var offset int64

			for !EOFError {
				fileData := make([]byte, FILE_DATA_ON_READ_PART_SIZE)

				bytesRead, err := file.ReadAt(fileData, offset)
				if err != nil {
					if errors.Is(err, io.EOF) {
						EOFError = true
					} else {
						log.Fatal(fmt.Errorf("compress: unable to read file bytes at %d. error: %w", offset, err))
					}
				}

				for i := range bytesRead {
					data = append(data, fileData[i])
				}

				offset += FILE_DATA_ON_READ_PART_SIZE
			}
			file.Close()
		}
		//[0 fileName([]byte) fileData([]byte) 0 fileName([]byte) fileData([]byte)]
		folderData = append(folderData, FILE_NAME_SEPARATOR_BYTE)   // next file data, byte-separator
		folderData = append(folderData, []byte(dirEntry.Name())...) // file name
		folderData = append(folderData, FILE_NAME_SEPARATOR_BYTE)   // separate name
		folderData = append(folderData, data...)                    // file data
	}

	compressedFolderData := folderData //TODO: byte compression

	//* create compressed file
	newFile, err := os.Create(fmt.Sprintf("%s/%s.%s", OUTPUT_FOLDER, folderName, FILE_EXTENSION))
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to create compressed file. error: %w", err))
	}

	//* write compressed bytes into new file
	_, err = newFile.Write(compressedFolderData)
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to write compressed data into created file. error: %w", err))
	}

	return compressedFolderData, time.Since(compressionStartTime).Milliseconds()
}
