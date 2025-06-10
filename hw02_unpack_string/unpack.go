package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func IsLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runes := []rune(str)
	currentPos := 0

	var builder strings.Builder

	for {
		currentRune := runes[currentPos]

		if unicode.IsDigit(currentRune) {
			return "", ErrInvalidString
		}

		if currentPos+1 == len(runes) {
			if IsLetter(currentRune) {
				builder.WriteString(string(currentRune))
			}

			return builder.String(), nil
		}

		nextRune := runes[currentPos+1]

		if unicode.IsDigit(nextRune) {
			digit, ok := strconv.Atoi(string(nextRune))

			if ok != nil {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(currentRune) {
				return "", ErrInvalidString
			}

			repeated := strings.Repeat(string(currentRune), digit)
			builder.WriteString(repeated)

			digitIsLast := currentPos+2 == len(runes)
			if digitIsLast {
				return builder.String(), nil
			}

			currentPos += 2

			continue
		}

		if IsLetter(currentRune) {
			builder.WriteString(string(currentRune))
		}

		currentPos++
	}
}
