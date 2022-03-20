package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/example/stringutil"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	const DoWrite, DoNothing, EmptyString = -1, 0, ""
	var builder strings.Builder
	actionNumber := DoWrite

	r := []rune(s)
	for i := len(r) - 1; i >= 0; i-- {
		if unicode.IsDigit(r[i]) {
			if (actionNumber >= DoNothing) || (i == 0) {
				return EmptyString, ErrInvalidString
			}
			actionNumber, _ = strconv.Atoi(string(r[i]))
		} else {
			switch actionNumber {
			case DoWrite:
				builder.WriteRune(r[i])
			case DoNothing:
				actionNumber = DoWrite
			default:
				builder.WriteString(strings.Repeat(string(r[i]), actionNumber))
				actionNumber = DoWrite
			}
		}
	}

	return stringutil.Reverse(builder.String()), nil
}
