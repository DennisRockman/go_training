package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	goldenMap := Environment{
		"BAR":   {Value: "bar", NeedRemove: false},
		"EMPTY": {Value: "", NeedRemove: false},
		"FOO":   {Value: "   foo" + "\n" + "with new line", NeedRemove: false},
		"HELLO": {Value: "\"hello\"", NeedRemove: false},
		"UNSET": {Value: "", NeedRemove: true},
	}
	envMap, err := ReadDir("./testdata/env")
	require.NoError(t, err)
	require.Equal(t, goldenMap, envMap)
}
