package styledconsole

import (
	"bufio"
	"io"
	"os"
)

type typedKey struct {
	KeyType   string
	Character rune
	ArrowKey  rune
}

func getKey() (*typedKey, error) {
	reader := bufio.NewReader(os.Stdin)

	// We do not handle the error yet but when calling ReadRune()
	peekedBytes, _ := reader.Peek(1)

	if len(peekedBytes) == 1 && peekedBytes[0] == '\033' && reader.Buffered() == 3 {
		peekedBytes, _ := reader.Peek(3)
		// We test for an escape sequence (byte 91 is "[")
		if len(peekedBytes) == 3 && peekedBytes[0] == '\033' && peekedBytes[1] == 91 {
			// We have an escape sequence, we discard it.
			_, _ = reader.Discard(3)

			// The escape codes for arrows are <esc>[A, <esc>[B, <esc>[C and <esc>[D
			switch peekedBytes[2] {
			case 'A':
				// advance reader by 3 bytes
				return &typedKey{KeyType: "arrowKey", ArrowKey: '↑'}, nil
			case 'B':
				return &typedKey{KeyType: "arrowKey", ArrowKey: '↓'}, nil
			case 'C':
				return &typedKey{KeyType: "arrowKey", ArrowKey: '→'}, nil
			case 'D':
				return &typedKey{KeyType: "arrowKey", ArrowKey: '←'}, nil
			}

			// The escape sequence is not handled yet.
			return nil, nil
		}
	}

	inputRune, _, err := reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return &typedKey{KeyType: "EOF"}, nil
		}

		return nil, err
	}

	return &typedKey{KeyType: "char", Character: inputRune}, nil
}
