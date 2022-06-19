package main

import (
	"flag"
	"os"
	"strings"
)

var (
	fileToWrite string
	envString   string
)

func init() {
	flag.StringVar(&fileToWrite, "file", "", "file to write to")
	flag.StringVar(&envString, "env_list", "", "environment variables list via ,")
}

func main() {
	flag.Parse()

	f, err := os.Create(fileToWrite)
	if err != nil {
		return
	}
	defer f.Close()

	envVariables := strings.Split(envString, ",")
	for _, envVar := range envVariables {
		_, err = f.WriteString(envVar + "=" + os.Getenv(envVar) + "\n")
	}
}
