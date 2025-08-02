package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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
// TODO: instead of using bytes slice for storing compressed data, write in output file during compression

// theres no actuall compression yet
func compressFileData(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s. err: %w", filePath, err)
	}
	defer file.Close()

	var EOFError bool
	var offset int64
	var fileDataLenUint uint32
	var fileData []byte

	for !EOFError {
		dataPart := make([]byte, FILE_DATA_ON_READ_PART_SIZE)

		bytesRead, err := file.ReadAt(dataPart, offset)
		if err != nil {
			if errors.Is(err, io.EOF) {
				EOFError = true
			} else {
				return nil, fmt.Errorf("compress: unable to read file bytes at %d. error: %w", offset, err)
			}
		}

		for i := 0; i < bytesRead; i++ {
			fileDataLenUint++
			fileData = append(fileData, dataPart[i])
		}

		offset += FILE_DATA_ON_READ_PART_SIZE
	}

	fileNameLenUint := uint16(len(filePath))
	fileNameLenBytes := make([]byte, FILE_NAME_BYTES_LENGTH)
	fileDataLenBytes := make([]byte, FILE_DATA_BYTES_LENGTH)

	binary.BigEndian.PutUint16(fileNameLenBytes, fileNameLenUint)
	binary.BigEndian.PutUint32(fileDataLenBytes, fileDataLenUint)

	var data []byte

	data = append(data, fileNameLenBytes...) // length of file path
	data = append(data, []byte(filePath)...) // file path
	data = append(data, fileDataLenBytes...) // length of file data
	data = append(data, fileData...)         // file data

	return data, nil
}

func compressFolderFilesData(folderPath string) ([]byte, error) {
	folderFiles, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read directiry: %s. err: %w", folderPath, err)
	}

	filesData := []byte{}

	for _, file := range folderFiles {
		filePath := fmt.Sprintf("%s/%s", folderPath, file.Name())

		var compressedData []byte

		if file.IsDir() {
			compressedData, err = compressFolderFilesData(filePath)
			if err != nil {
				return nil, err
			}
		} else {
			compressedData, err = compressFileData(filePath)
			if err != nil {
				return nil, err
			}
		}

		filesData = append(filesData, compressedData...)
	}

	return filesData, nil
}

// return function work time in milliseconds
func compressWRTAIO(folderName string) int64 {
	//* get folder files
	inputFiles, err := os.ReadDir(folderName)
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to read directory. error: %w", err))
	}

	compressionStartTime := time.Now()

	compressedFilesData := []byte{}

	for _, dirEntry := range inputFiles {
		filePath := fmt.Sprintf("%s/%s", folderName, dirEntry.Name())

		var data []byte
		var fileDataLenUint uint32

		if dirEntry.IsDir() {
			data, err = compressFolderFilesData(filePath)
			if err != nil {
				log.Printf("skipped file: %s. error: %v", filePath, err)
				continue
			}
		} else {
			//* file open
			data, err = compressFileData(filePath)
			if err != nil {
				log.Printf("skipped file: %s. error: %v", filePath, err)
				continue
			}
		}

		filePathLenUint := uint16(len(filePath))
		fileDataLenUint = uint32(len(data))

		filePathLenBytes := make([]byte, FILE_NAME_BYTES_LENGTH)
		fileDataLenBytes := make([]byte, FILE_DATA_BYTES_LENGTH)

		binary.BigEndian.PutUint16(filePathLenBytes, filePathLenUint)
		binary.BigEndian.PutUint32(fileDataLenBytes, fileDataLenUint)

		compressedFilesData = append(compressedFilesData, data...) // compressed file data
	}

	//* create compressed file
	newFile, err := os.Create(fmt.Sprintf("%s/%s.%s", OUTPUT_FOLDER, folderName, FILE_EXTENSION))
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to create compressed file. error: %w", err))
	}

	//* write compressed bytes into new file
	_, err = newFile.Write(compressedFilesData)
	if err != nil {
		log.Fatal(fmt.Errorf("compress: unable to write compressed data into created file. error: %w", err))
	}

	return time.Since(compressionStartTime).Milliseconds()
}

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

	for decompressedFilePath, bytes := range decompData {
		s := []string{OUTPUT_FOLDER}
		s = append(s, strings.Split(decompressedFilePath, "/")...)

		for i := 0; i < len(s)-1; i++ {
			path := strings.Join(s[:i+1], "/")
			if _, err := os.ReadDir(path); err != nil {
				if err := os.Mkdir(path, 0644); err != nil {
					log.Fatal(fmt.Errorf("decompress: unable to create folder from compressed file: %s. error:  %w", path, err))
				}
			}
		}

		file, err := os.Create(fmt.Sprintf("%s/%s", OUTPUT_FOLDER, decompressedFilePath)) // output/innerFolder/file.txt
		log.Println(decompressedFilePath)
		if err != nil {
			log.Fatal(fmt.Errorf("decompress: unable to create files with decompressed data. error: %w", err))
		}
		defer file.Close()

		file.Write(bytes)
	}

	return time.Since(decompressionStartTime).Milliseconds()
}
