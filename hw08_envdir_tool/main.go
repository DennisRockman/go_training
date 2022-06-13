package main

import (
	"flag"
	"os"
)

const Equal = "="

func main() {
	flag.Parse()
	args := flag.Args()
	envMap, err := ReadDir(args[0])
	if err != nil {
		os.Exit(4)
	}
	returnCode := RunCmd(args[1:], envMap)
	if returnCode > 0 {
		os.Exit(returnCode)
	}
}
