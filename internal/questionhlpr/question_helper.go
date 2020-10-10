package questionhlpr

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/corentindeboisset/styledconsole/internal/style"
	"golang.org/x/crypto/ssh/terminal"
)

type Question struct {
	Label         string
	IsClosed      bool
	IsHidden      bool
	Answers       []string
	DefaultAnswer string
	Validator     func(string) bool
}

var greenStyle, yellowStyle, redStyle *style.OutputStyle

func init() {
	greenStyle = style.NewOutputStyle("fg=green")
	yellowStyle = style.NewOutputStyle("fg=yellow")
	redStyle = style.NewOutputStyle("fg=red")
}

func AskQuestion(q Question) (string, error) {
	if q.IsClosed && len(q.Answers) > 1 {
		ret, err := askClosedQuestion(q)
		if err != nil {
			return "", err
		}

		return ret, nil
	} else if !q.IsClosed {
		if q.IsHidden {
			for {
				ret, err := askHiddenQuestion(q)
				if err != nil {
					return "", err
				}
				if q.Validator == nil || q.Validator(ret) {
					return ret, nil
				} else {
					fmt.Printf("%s\n", redStyle.Apply("This answer is invalid."))
				}
			}
		} else {
			for {
				ret, err := askRegularQuestion(q)
				if err != nil {
					return "", err
				}
				if q.Validator == nil || q.Validator(ret) {
					return ret, nil
				} else {
					fmt.Printf("%s\n", redStyle.Apply("This answer is invalid."))
				}
			}
		}
	}

	// Invalid question, empty answer.
	// This results from an error of implementation and should never happen
	return "", errors.New("The question object is invalid")
}

func askClosedQuestion(q Question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open a prompt outside of a TTY")
	}

	fmt.Printf(
		"\n%s [%s]",
		greenStyle.Apply(strings.TrimSpace(q.Label)),
		yellowStyle.Apply("default"),
	)

	return q.Answers[0], nil
}

func askHiddenQuestion(q Question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open a prompt outside of a terminal")
	}

	fmt.Printf("\n%s :\n > ", greenStyle.Apply(strings.TrimSpace(q.Label)))
	answerBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	// The typed line break is hidden so we have to force it
	fmt.Print("\n")

	if err != nil {
		return "", fmt.Errorf("There was an error reading the stdin (%s)", err)
	}

	return string(answerBytes), nil
}

func askRegularQuestion(q Question) (string, error) {
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return "", errors.New("Cannot open a prompt outside of a terminal")
	}

	fmt.Printf("\n%s :\n > ", greenStyle.Apply(strings.TrimSpace(q.Label)))

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	// A line break will always be typed so no need to force it

	if err != nil {
		return "", fmt.Errorf("There was an error reading the stdin (%s)", err)
	}

	return answer[:len(answer)-1], nil
}
