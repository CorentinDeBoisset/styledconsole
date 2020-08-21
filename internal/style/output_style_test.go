package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApplyStyle checks that we can apply a test on a string
func TestApplyStyle(t *testing.T) {
	mystyle := OutputStyle{Foreground: "green", Background: "red", Options: []string{"bold"}}
	assert.Equal(
		t,
		"\033[1;32;41mThis is a text.\033[22;39;49m",
		mystyle.Apply("This is a text."),
	)
}

// TestApplyInvalidStyle checks that invalid style properties will not be applied
func TestApplyInvalidStyle(t *testing.T) {
	mystyle := OutputStyle{Foreground: "blau", Background: "rojo", Options: []string{"gras"}}
	assert.Equal(
		t,
		"This is a text.",
		mystyle.Apply("This is a text."),
	)

	mystyle = OutputStyle{Foreground: "blue", Background: "rojo", Options: []string{"gras"}}
	assert.Equal(
		t,
		"\033[34mThis is a text.\033[39m",
		mystyle.Apply("This is a text."),
	)
}

// TestMergeStyles checks we can merge two styles together
func TestMergeStyles(t *testing.T) {
	assert := assert.New(t)
	firstStyle := OutputStyle{Foreground: "blue", Background: "blue", Options: []string{"bold", "underscore", "blink"}}
	secondStyle := OutputStyle{Foreground: "red", Background: "red", Options: []string{"reverse", "conceal", "underscore"}}

	secondStyle.MergeBase(firstStyle)
	assert.Equal(
		OutputStyle{Foreground: "red", Background: "red", Options: []string{"reverse", "conceal", "underscore", "bold", "blink"}},
		secondStyle,
	)

	fgStyle := OutputStyle{Foreground: "green"}
	fgStyle.MergeBase(firstStyle)
	assert.Equal(
		OutputStyle{Foreground: "green", Background: "blue", Options: []string{"bold", "underscore", "blink"}},
		fgStyle,
	)

	bgStyle := OutputStyle{Background: "green"}
	bgStyle.MergeBase(firstStyle)
	assert.Equal(
		OutputStyle{Foreground: "blue", Background: "green", Options: []string{"bold", "underscore", "blink"}},
		bgStyle,
	)
}
