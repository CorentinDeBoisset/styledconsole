package styledprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew checks we can build a style from a string
func TestNew(t *testing.T) {
	assert := assert.New(t)

	// Test a simple, valid style
	assert.Equal(&OutputStyle{foreground: "red", background: "green", handleHrefGracefully: true}, NewOutputStyle("bg=green;fg=red"))

	// Test a style with options and href
	assert.Equal(
		&OutputStyle{foreground: "ieua", background: "aie", href: "http://github.com", handleHrefGracefully: true, options: map[string]bool{"bold": true, "italic": true}},
		NewOutputStyle("bg=aie;fg=ieua;href=http://github.com;options=bold,italic"),
	)

	// Test an invalid style
	assert.Equal((*OutputStyle)(nil), NewOutputStyle("toto=titi;fg=red"))
}

// TestApplyStyle checks that we can apply a test on a string
func TestApplyStyle(t *testing.T) {
	mystyle := OutputStyle{foreground: "green", background: "red", options: map[string]bool{"bold": true}}
	assert.Equal(
		t,
		"\033[1;32;41mThis is a text.\033[22;39;49m",
		mystyle.Apply("This is a text."),
	)
}

// TestApplyInvalidStyle checks that invalid style properties will not be applied
func TestApplyInvalidStyle(t *testing.T) {
	mystyle := OutputStyle{foreground: "blau", background: "rojo", options: map[string]bool{"gras": true}}
	assert.Equal(
		t,
		"This is a text.",
		mystyle.Apply("This is a text."),
	)

	mystyle = OutputStyle{foreground: "blue", background: "rojo", options: map[string]bool{"gras": true}}
	assert.Equal(
		t,
		"\033[34mThis is a text.\033[39m",
		mystyle.Apply("This is a text."),
	)
}

// TestMergeStyles checks we can merge two styles together
func TestMergeStyles(t *testing.T) {
	assert := assert.New(t)
	firstStyle := "fg=blue;bg=blue;options=bold,underscore,blink"
	secondStyle := "fg=red;bg=red;options=reverse,conceal,underscore"

	assert.Equal(
		&OutputStyle{
			foreground:           "red",
			background:           "red",
			handleHrefGracefully: true,
			options:              map[string]bool{"reverse": true, "conceal": true, "underscore": true, "bold": true, "blink": true},
		},
		NewOutputStyle(firstStyle+";"+secondStyle),
	)

	fgStyle := "fg=green"
	assert.Equal(
		&OutputStyle{
			foreground:           "green",
			background:           "blue",
			handleHrefGracefully: true,
			options:              map[string]bool{"bold": true, "underscore": true, "blink": true},
		},
		NewOutputStyle(firstStyle+";"+fgStyle),
	)

	bgStyle := "bg=green"
	assert.Equal(
		&OutputStyle{
			foreground:           "blue",
			background:           "green",
			handleHrefGracefully: true,
			options:              map[string]bool{"bold": true, "underscore": true, "blink": true},
		},
		NewOutputStyle(firstStyle+";"+bgStyle),
	)
}
