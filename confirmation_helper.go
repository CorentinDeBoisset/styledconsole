package styledconsole

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

func askConfirm(label string, defaultAnswer *bool) (bool, error) {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return false, errors.New("cannot open a prompt outside of a terminal")
	}

	for {
		var options string
		if defaultAnswer != nil {
			if *defaultAnswer {
				options = "Y/n"
			} else {
				options = "y/N"
			}
		} else {
			options = "y/n"
		}
		fmt.Printf("%s [%s]: ", greenStyle.Apply(strings.TrimSpace(label)), yellowStyle.Apply(options))

		reader := bufio.NewReader(os.Stdin)
		textAnswer, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF && defaultAnswer != nil {
				return *defaultAnswer, nil
			}
			return false, fmt.Errorf("there was an error reading the stdin: %w", err)
		}

		if strings.ToLower(textAnswer) == "yes\n" || strings.ToLower(textAnswer) == "y\n" {
			return true, nil
		} else if strings.ToLower(textAnswer) == "no\n" || strings.ToLower(textAnswer) == "n\n" {
			return false, nil
		} else if textAnswer == "\n" && defaultAnswer != nil {
			return *defaultAnswer, nil
		}
	}
}
