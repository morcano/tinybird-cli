package utils

import (
	"bufio"
	"fmt"
	"os"
)

type DecoderFunc func(ciphertext string, key string) (string, error)

// ValidateFile checks if the file exists at the given path
func ValidateFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("the file %s does not exist", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("unable to open the file %s: %v", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 1
	for scanner.Scan() {
		token := scanner.Text()
		if token == "" {
			return fmt.Errorf("empty token found on line %d", lineNumber)
		}
		lineNumber++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	if lineNumber == 1 {
		return fmt.Errorf("the file %s is empty", filePath)
	}

	return nil
}

func GetItems(filePath, decryptKey string, decoder DecoderFunc) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var items []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		item := scanner.Text()
		if decryptKey != "" && decoder != nil {
			item, err = decoder(item, decryptKey)
			if err != nil {
				fmt.Printf("Error decoding item: %v\n", err)
				continue // Optionally, decide whether to continue or return on decryption error
			}
		}
		items = append(items, item)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return items, nil
}
