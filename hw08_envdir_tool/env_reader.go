package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var ErrInvalidFilename = errors.New("invalid filename")

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envMap := make(Environment)
	for _, fileInfo := range filesInfo {
		fileName := fileInfo.Name()
		if strings.Contains(fileName, Equal) {
			return nil, ErrInvalidFilename
		}
		filePath := filepath.Join(dir, fileName)
		err := readFile(fileName, filePath, envMap)
		if err != nil {
			return nil, err
		}
	}

	return envMap, nil
}

func readFile(fileName string, filePath string, environment Environment) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan() // read first line
	if err := scanner.Err(); err != nil {
		return err
	}

	clearBytes := bytes.ReplaceAll(scanner.Bytes(), []byte("\x00"), []byte("\n"))
	clearText := strings.TrimRight(string(clearBytes), " ")
	environment[fileName] = EnvValue{
		Value:      clearText,
		NeedRemove: isEmpty(fileStat.Size()),
	}

	return nil
}

func isEmpty(fileSize int64) bool {
	return fileSize <= 0
}
