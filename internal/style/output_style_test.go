package style

import (
	"fmt"
	"strconv"
	"testing"
)

func assertEqualsString(t *testing.T, a string, b string) {
	if a == b {
		return
	}

	t.Error(fmt.Sprintf("Assertion failed:\n  %v != %v", strconv.QuoteToASCII(a), strconv.QuoteToASCII(b)))
}

// TestApplyStyle checks that we can apply a test on a string
func TestApplyStyle(t *testing .T) {
	mystyle := OutputStyle{Foreground: "green", Background: "red", Options: []string{"bold"}}

	assertEqualsString(t, mystyle.Apply("This is a text."), "\033[1;32;41mThis is a text.\033[22;39;49m")
}
