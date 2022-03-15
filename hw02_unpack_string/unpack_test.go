package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: " 5", expected: "     "},
		{input: "a1a😀5", expected: "aa😀😀😀😀😀"},
		{input: "a1", expected: "a"},
		{input: "a-1", expected: "a-"},
		{input: "a0.5", expected: "....."},
		{input: "a-0.5", expected: "a....."},
		{input: "a0", expected: ""},
		{input: "a0b0c0", expected: ""},
		{input: "a1-0", expected: "a"},
		{input: "дЯ1Дя1 ФЁ5Д0р", expected: "дЯДя ФЁЁЁЁЁр"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackTwiceCall(t *testing.T) {
	tests := []struct {
		input1  string
		input2  string
		isEqual bool
	}{
		{input1: " 1 2 3", input2: " 6", isEqual: true},
		{input1: " 0 0 0", input2: " 0", isEqual: true},
		{input1: "$1$2$3", input2: "$3$2$1", isEqual: true},
		{input1: "aaa", input2: "a3", isEqual: true},
		{input1: "", input2: "a0", isEqual: true},
		{input1: "-1-2-3", input2: "-6", isEqual: true},
		{input1: "ABBA", input2: "A1B2AA0", isEqual: true},
		{input1: "у3ор", input2: "уу2ор", isEqual: true},
		{input1: " 0", input2: " 1-1", isEqual: false},
		{input1: "$1 $2 $3", input2: "$3 $2 $1", isEqual: false},
		{input1: "A1", input2: "А1", isEqual: false}, // English and Russian
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input1, func(t *testing.T) {
			result1, err1 := Unpack(tc.input1)
			result2, err2 := Unpack(tc.input2)
			require.NoError(t, err1)
			require.NoError(t, err2)
			if tc.isEqual {
				require.Equal(t, result1, result2)
			} else {
				require.NotEqual(t, result1, result2)
			}
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "5", "00"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
