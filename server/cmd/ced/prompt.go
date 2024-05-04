package main

import (
	"bufio"
	"io"
	"unicode"
)

func askYesNo(prompt string, read io.Reader, write io.Writer) (bool, error) {
	reader := bufio.NewReader(read)

	for {
		if _, err := write.Write([]byte(prompt + "\n")); err != nil {
			return false, err
		}

		rune, _, err := reader.ReadRune()
		if err != nil {
			return false, err
		}

		rune = unicode.ToLower(rune)
		if rune == 'y' {
			return true, nil
		}
		if rune == 'n' {
			return false, nil
		}
		continue
	}
}
