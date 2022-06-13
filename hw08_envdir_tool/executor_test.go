package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	separator   = ","
	utilPathArg = "./executor_test_util/executor_test_util"
	fileName    = "env_data.txt"
	filePathArg = "-file=" + fileName
	envListArg  = "-env_list="
)

func TestRunCmd(t *testing.T) {
	envMap := Environment{
		"FIRST":   {Value: "1", NeedRemove: false},
		"SECOND":  {Value: "2", NeedRemove: false},
		"_EMPTY_": {Value: "", NeedRemove: false},
		"DELETE":  {Value: "", NeedRemove: true},
	}

	var sb strings.Builder
	sb.WriteString(envListArg)
	for key := range envMap {
		sb.WriteString(key)
		sb.WriteString(separator)
	}
	envList := strings.TrimRight(sb.String(), separator)

	cmd := []string{utilPathArg, filePathArg, envList}
	returnCode := RunCmd(cmd, envMap)
	require.Zero(t, returnCode)

	envBytes, _ := os.ReadFile(fileName)
	envString := string(envBytes)
	defer os.Remove(fileName)

	for key, value := range envMap {
		sb.Reset()
		sb.WriteString(key)
		sb.WriteString(Equal)
		sb.WriteString(value.Value)
		require.True(t, strings.Contains(envString, sb.String()))
	}
}
