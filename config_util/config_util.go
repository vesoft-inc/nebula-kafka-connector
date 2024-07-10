package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

func WriteToFile(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

func DecompressData(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(data)
	gz, err := gzip.NewReader(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gz.Close()

	var decompressedData bytes.Buffer
	_, err = decompressedData.ReadFrom(gz)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	return decompressedData.Bytes(), nil
}

func CompressData(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var compressedData bytes.Buffer
	gz := gzip.NewWriter(&compressedData)
	_, err = io.Copy(gz, file)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}
	gz.Close()

	return compressedData.Bytes(), nil
}

func Serialize(stream []byte) (result string) {
	result = ""
	for _, b := range stream {
		result += fmt.Sprintf("%d,", b)
	}
	result = strings.TrimSuffix(result, ",")
	return result
}

func CompareFiles(file1Path, file2Path string) (bool, error) {
	file1, err := os.Open(file1Path)
	if err != nil {
		return false, fmt.Errorf("failed to open file1: %v", err)
	}
	defer file1.Close()

	file2, err := os.Open(file2Path)
	if err != nil {
		return false, fmt.Errorf("failed to open file2: %v", err)
	}
	defer file2.Close()

	file1Stat, err := file1.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file1 info: %v", err)
	}

	file2Stat, err := file2.Stat()
	if err != nil {
		return false, fmt.Errorf("failed to get file2 info: %v", err)
	}

	if file1Stat.Size() != file2Stat.Size() {
		return false, nil
	}

	file1Data := make([]byte, file1Stat.Size())
	_, err = file1.Read(file1Data)
	if err != nil {
		return false, fmt.Errorf("failed to read file1 data: %v", err)
	}

	file2Data := make([]byte, file2Stat.Size())
	_, err = file2.Read(file2Data)
	if err != nil {
		return false, fmt.Errorf("failed to read file2 data: %v", err)
	}

	return bytes.Equal(file1Data, file2Data), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path as a command line argument.")
		return
	}

	filePath := os.Args[1]
	compressedData, err := CompressData(filePath)
	if err != nil {
		fmt.Printf("Failed to compress data: %v\n", err)
		return
	}

	err = WriteToFile("encoded.txt", []byte(Serialize(compressedData)))
	if err != nil {
		fmt.Printf("Failed to write compressed data to file: %v\n", err)
	}

	decompressedData, err := DecompressData(compressedData)
	if err != nil {
		fmt.Printf("Failed to decompress data: %v\n", err)
		return
	}

	err = WriteToFile("decoded.txt", decompressedData)
	if err != nil {
		fmt.Printf("Failed to write decompressed data to file: %v\n", err)
	}

	areEqual, err := CompareFiles(filePath, "decoded.txt")
	if err != nil || areEqual == false{
		fmt.Printf("Files are not equal or fail to compare: %v\n", err)
		return
	}

	return
}
