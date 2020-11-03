package styledconsole

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type question struct {
	Label         string
	IsClosed      bool
	IsHidden      bool
	Choices       []string
	DefaultChoice int
	DefaultAnswer string
	Validator     func(string) bool
}

func askQuestion(q question) (string, error) {
	if q.IsClosed && len(q.Choices) > 1 {
		ret, err := askClosedQuestion(q)
		if err != nil {
			return "", err
		}

		return ret, nil
	} else if !q.IsClosed {
		var ret string
		var err error
		for {
			if q.IsHidden {
				ret, err = askHiddenQuestion(q)
			} else {
				ret, err = askRegularQuestion(q)
			}

			if err != nil {
				if err == io.EOF && (q.Validator == nil || q.Validator(ret)) {
					// Handle empty buffer gracefully if possible
					return ret, nil
				}

				return "", err
			}

			if q.Validator == nil || q.Validator(ret) {
				return ret, nil
			} else {
				fmt.Printf("%s\n", redStyle.Apply("This answer is invalid."))
			}
		}
	}

	// Invalid question, empty answer.
	// This results from an error of implementation in the library and should never happen
	return "", errors.New("The question object is invalid")
}

func askClosedQuestion(q question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open an interacive prompt outside of a TTY")
	}

	width, height := getWinsize()

	oldState, err := terminal.MakeRaw(int(os.Stdout.Fd()))
	if err != nil {
		return "", fmt.Errorf("There was an error switching terminal to raw mode (%s)", err)
	}

	if height < 3 || width < 20 {
		for {
			fmt.Print("Terminal is too small... Resize and press a key.\n")
			key, err := getKey()
			if err != nil {
				return "", fmt.Errorf("There was an error reading a character from stdin (%s)", err)
			}

			if key.KeyType == "EOF" {
				// We skip and move on to the next step
				break
			}

			width, height = getWinsize()
		}
	}

	// Prepare the list of printable options
	printableChoices := []string{}
	for _, choice := range q.Choices {
		if len(choice) < width-3 {
			printableChoices = append(printableChoices, choice)
		} else {
			printableChoices = append(printableChoices, fmt.Sprintf("%s…", choice[:width-5]))
		}
	}

	// Run the display loop
	selectedIndex := -1
	choiceCount := len(printableChoices)
	scrollWindowHeight := getScrollWindowHeight(choiceCount, height)
	highlightedIndex := 0
	scroll := 0

	if q.DefaultChoice >= 0 && q.DefaultChoice < choiceCount-1 {
		// Set the default answer, and adapt the initial scrolling
		highlightedIndex = q.DefaultChoice
		if highlightedIndex > scrollWindowHeight {
			if highlightedIndex < choiceCount-2 {
				scroll = highlightedIndex - scrollWindowHeight
			} else {
				scroll = choiceCount - scrollWindowHeight - 2
			}
		}
	}

	hideCursor()
	for selectedIndex == -1 {
		clearWindowFromCursor()
		fmt.Printf("%s:", greenStyle.Apply(q.Label))

		// Print the first line, either the first choice or a "↑"
		if scroll > 0 {
			fmt.Print("\n\033[1000D   ↑")
		} else {
			fmt.Print(formatClosedQuestionChoice(printableChoices[0], highlightedIndex == 0))
		}

		// Print some choices
		for i := scroll + 1; i <= scroll+scrollWindowHeight; i++ {
			fmt.Print(formatClosedQuestionChoice(printableChoices[i], highlightedIndex == i))
		}

		// Print the last line, either the last choice or a "↓"
		if scroll < choiceCount-scrollWindowHeight-2 {
			fmt.Print("\n\033[1000D   ↓")
		} else {
			fmt.Print(formatClosedQuestionChoice(printableChoices[choiceCount-1], highlightedIndex == choiceCount-1))
		}

		// Put the cursor back at the beginning
		fmt.Printf("\033[%dA\033[1000D", scrollWindowHeight+2)

		for {
			typedKey, err := getKey()

			// Re-parse the height in case the user resized their terminal
			_, height = getWinsize()
			scrollWindowHeight = getScrollWindowHeight(choiceCount, height)

			if err != nil || typedKey == nil {
				return "", fmt.Errorf("There was an error reading user input (%s)", err)
			}
			if typedKey.KeyType == "EOF" {
				if q.DefaultChoice >= 0 && q.DefaultChoice < choiceCount-1 {
					return q.Choices[q.DefaultChoice], nil
				}

				return "", errors.New("Error parsing user activity from Stdin (EOF)")
			}

			if typedKey.KeyType == "arrowKey" && typedKey.ArrowKey == '↑' {
				// Up
				if highlightedIndex == 0 {
					highlightedIndex = choiceCount - 1
					scroll = choiceCount - scrollWindowHeight - 2 // scroll to the bottom
				} else {
					highlightedIndex -= 1
					// Update scrolling if necessary
					if scroll >= highlightedIndex {
						if highlightedIndex > 1 {
							scroll = highlightedIndex - 1
						} else {
							scroll = 0
						}
					}
				}
				break
			} else if typedKey.KeyType == "arrowKey" && typedKey.ArrowKey == '↓' {
				// Down
				if highlightedIndex == choiceCount-1 {
					highlightedIndex = 0
					scroll = 0 // scroll to the top
				} else {
					highlightedIndex += 1
					// Update scrolling if necessary
					if scroll <= highlightedIndex-scrollWindowHeight-1 {
						if highlightedIndex < choiceCount-2 {
							scroll = highlightedIndex - scrollWindowHeight
						} else {
							scroll = choiceCount - scrollWindowHeight - 2
						}
					}
				}
				break
			} else if typedKey.KeyType == "char" && (typedKey.Character == '\r' || typedKey.Character == '\n' || typedKey.Character == ' ') {
				selectedIndex = highlightedIndex
				break
			} else if typedKey.KeyType == "char" && typedKey.Character == 3 {
				// Ctrl-C
				showCursor()
				fmt.Printf("\033[%dB\033[1000D", scrollWindowHeight+3)
				_ = terminal.Restore(int(os.Stdout.Fd()), oldState)
				_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}
		}
	}
	fmt.Printf("\033[%dB\033[1000D", scrollWindowHeight+3)
	showCursor()
	_ = terminal.Restore(int(os.Stdout.Fd()), oldState)

	return q.Choices[selectedIndex], nil
}

func askHiddenQuestion(q question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open a prompt outside of a terminal")
	}

	fmt.Printf("\n%s :\n > ", greenStyle.Apply(strings.TrimSpace(q.Label)))
	answerBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	// The typed line break is hidden so we have to force it
	fmt.Print("\n")

	if err != nil {
		if err == io.EOF {
			return string(answerBytes), err
		}

		return "", fmt.Errorf("There was an error reading the stdin (%s)", err)
	}

	return string(answerBytes), nil
}

func askRegularQuestion(q question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open a prompt outside of a terminal")
	}

	var prompt string
	if q.DefaultAnswer != "" {
		prompt = fmt.Sprintf("\n%s [%s]:\n > ", greenStyle.Apply(strings.TrimSpace(q.Label)), yellowStyle.Apply(q.DefaultAnswer))
	} else {
		prompt = fmt.Sprintf("\n%s :\n > ", greenStyle.Apply(strings.TrimSpace(q.Label)))
	}
	fmt.Print(prompt)

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			return answer, err
		}

		return answer, fmt.Errorf("There was an error reading the stdin (%s)", err)
	}

	if answer == "\n" && q.DefaultAnswer != "" {
		return q.DefaultAnswer, nil
	}

	// Remove the ending "\n" before returning
	return answer[:len(answer)-1], nil
}

func getScrollWindowHeight(choiceCount int, termHeight int) int {
	if choiceCount > 12 && termHeight >= 13 {
		return 10
	} else if choiceCount+1 > termHeight {
		return termHeight - 3
	}

	return choiceCount - 2
}

func formatClosedQuestionChoice(label string, highlighted bool) string {
	if highlighted {
		return fmt.Sprintf("\n\033[1000D > %s", highlightedChoiceStyle.Apply(label))
	}

	return fmt.Sprintf("\n\033[1000D   %s", label)
}
