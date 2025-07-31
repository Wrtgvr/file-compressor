package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

/* so, I realized that byte-separators I used is SHIT.
2 bytes: length of file name,
N bytes: fileName,
4 bytes: length of file data,
M bytes: fileData,
*/

// maybe I need to separate some blocks of the code into different functions ://
// TODO: code refactoring
// ? this TODO and similar TODOs in other files will be here forever so every time I open this file it remind me of clearing code

func decompressWRTAIO(filepath string) int64 {
	//* open file

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(fmt.Errorf("decompress: unable to open file. does file exists? error: %w", err))
	}
	defer file.Close()

	decompressionStartTime := time.Now()

	//* decompression
	decompData := map[string][]byte{}

	var EOFError bool
	var offset int64

	for !EOFError {
		//* get file name
		// get file name length

		var fileNameLenUint uint16
		fileNameLenBytes := make([]byte, FILE_NAME_BYTES_LENGTH)

		_, err = file.ReadAt(fileNameLenBytes, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
				break
			}
			log.Fatal(fmt.Errorf("decompress: unable to get file name length bytes: %w", err))
		}
		offset += int64(FILE_NAME_BYTES_LENGTH)

		fileNameLenUint = binary.BigEndian.Uint16(fileNameLenBytes)

		// get file name
		var fileName string
		fileNameBytes := make([]byte, fileNameLenUint)

		_, err = file.ReadAt(fileNameBytes, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
				break
			}
			log.Fatal(fmt.Errorf("decompress: unable to get file name: %w", err))
		}
		offset += int64(fileNameLenUint)

		fileName = string(fileNameBytes)

		//* get file data
		// get file data length

		var fileDataLenUint uint32
		fileDataLenBytes := make([]byte, FILE_DATA_BYTES_LENGTH)

		_, err = file.ReadAt(fileDataLenBytes, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
				break
			}
			log.Fatal(fmt.Errorf("decompress: unable to get file data length bytes: %w", err))
		}
		offset += int64(FILE_DATA_BYTES_LENGTH)

		fileDataLenUint = binary.BigEndian.Uint32(fileDataLenBytes)

		// get file data
		fileData := make([]byte, fileDataLenUint)

		_, err := file.ReadAt(fileData, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
				break
			}
			log.Fatal(fmt.Errorf("decompress: unable to get file data: %w", err))
		}
		offset += int64(fileDataLenUint)

		decompData[fileName] = fileData
		fmt.Printf("File: %s, Bytes: %d\n", fileName, len(fileData))
	}

	for n := range decompData {
		fmt.Println(n) //!!
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
		var fileDataLenUint uint32

		if dirEntry.IsDir() {
			// if file is a folder then compress files from inner folder
			// btw is doesn't mean to compress data like "yo theres a folder"
			// so after decompression there won't be any folders which were in input folder
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
					fileDataLenUint++
					data = append(data, fileData[i])
				}

				offset += FILE_DATA_ON_READ_PART_SIZE
			}
			file.Close()
		}
		//[0 fileName([]byte) fileData([]byte) 0 fileName([]byte) fileData([]byte)]

		fileNameLenUint := uint16(len(dirEntry.Name()))
		fileNameLenBytes := make([]byte, FILE_NAME_BYTES_LENGTH)
		fileDataLenBytes := make([]byte, FILE_DATA_BYTES_LENGTH)

		binary.BigEndian.PutUint16(fileNameLenBytes, fileNameLenUint)
		binary.BigEndian.PutUint32(fileDataLenBytes, fileDataLenUint)

		folderData = append(folderData, fileNameLenBytes...)        // length of file name
		folderData = append(folderData, []byte(dirEntry.Name())...) // file name
		folderData = append(folderData, fileDataLenBytes...)        // length of file data
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
