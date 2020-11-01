package styledprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSubstring checks we can extract a substring safely
func TestGetSubstring(t *testing.T) {
	assert := assert.New(t)
	assert.Equal("aaa", getSubstring("zaaaz", 1, 4))
	assert.Equal("aaaz", getSubstring("zaaaz", 1, 50))
	assert.Equal("", getSubstring("zaaaz", 50, 51))
}

// TestEscapeTrailingBackslash checks we can remove trailing "\"" from texts
func TestEscapeTrailingBackslash(t *testing.T) {
	assert := assert.New(t)

	assert.Equal("noop", escapeTrailingBackslash("noop"))
	assert.Equal("super super\x00\x00", escapeTrailingBackslash("super super\\\\"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00", escapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00\x00\x00", escapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00\\\\"))
}

// TestFormatWithStyleWithoutTextBefore checks there are no errors of newline
func TestFormatWithStyleWithoutTextBefore(t *testing.T) {
	assert := assert.New(t)
	width := 20
	stack := newOutputStyleStack("")

	// Test line-cutting
	output := []string{}
	lastLineLength := 0
	addStringWithStyle("abc", width, &output, &lastLineLength, stack)
	assert.Equal(3, lastLineLength)
	assert.Equal([]string{"abc"}, output)

	output = []string{}
	lastLineLength = 0
	addStringWithStyle("supertoto\nabc\n", width, &output, &lastLineLength, stack)
	assert.Equal(0, lastLineLength)
	assert.Equal([]string{"supertoto           ", "abc                 ", ""}, output)

	output = []string{}
	lastLineLength = 0
	addStringWithStyle("abc\n\n\n", width, &output, &lastLineLength, stack)
	assert.Equal(0, lastLineLength)
	assert.Equal([]string{"abc                 ", ""}, output)

	output = []string{"iiiii"}
	lastLineLength = 5
	addStringWithStyle("super super super super super super super super super super", width, &output, &lastLineLength, stack)
	assert.Equal(4, lastLineLength)
	assert.Equal(
		[]string{"iiiiisuper super sup", "er super super super", " super super super s", "uper"},
		output,
	)

	// Test with style
	stack.Push("bg=green;fg=red")
	output = []string{"iiiii"}
	lastLineLength = 5
	addStringWithStyle("super super super super super super", width, &output, &lastLineLength, stack)
	assert.Equal(0, lastLineLength)
	assert.Equal(
		[]string{"iiiii\x1b[31;42msuper super sup\x1b[39;49m", "\x1b[31;42mer super super super\x1b[39;49m", ""},
		output,
	)
	stack.PopCurrent()

	// Test edge-cases
	output = []string{""}
	lastLineLength = 0
	addStringWithStyle(" ", width, &output, &lastLineLength, stack)
	assert.Equal(1, lastLineLength)
	assert.Equal([]string{" "}, output)

	output = []string{"abcdeabcde"}
	lastLineLength = 10
	addStringWithStyle("", width, &output, &lastLineLength, stack)
	assert.Equal(10, lastLineLength)
	assert.Equal([]string{"abcdeabcde"}, output)

	output = []string{"super super super super super super super super super long line"}
	lastLineLength = 63
	addStringWithStyle("abc", width, &output, &lastLineLength, stack)
	assert.Equal(3, lastLineLength)
	assert.Equal([]string{"super super super super super super super super super long line", "abc"}, output)

	output = []string{"iii"}
	lastLineLength = -10
	addStringWithStyle("abc", width, &output, &lastLineLength, stack)
	assert.Equal(6, lastLineLength)
	assert.Equal([]string{"iiiabc"}, output)
}

// TestFormatTextWithoutDefault checks we can render a full text using style tags
func TestFormatTextWithoutDefault(t *testing.T) {
	assert := assert.New(t)
	width := 20

	assert.Equal(
		[]string{"great text"},
		formatText("great text", width, ""),
	)
	assert.Equal(
		[]string{"awesome text        ", "on                  ", "multiple lines."},
		formatText("awesome text\non\nmultiple lines.", width, ""),
	)
	assert.Equal(
		[]string{"\x1b[31mawesome text\x1b[39m"},
		formatText("<fg=red>awesome text</>", width, ""),
	)
	assert.Equal(
		[]string{
			"Some                          ",
			"text                          ",
			"that can handle \x1b[31mmulti-line    \x1b[39m",
			"\x1b[31mstyling.\x1b[39m This is a very long l",
			"ine.",
		},
		formatText("Some\ntext\nthat can handle <fg=red>multi-line\nstyling.</> This is a very long line.", 30, ""),
	)
	assert.Equal(
		[]string{
			"awesome text \x1b[31mwith \x1b[39m\x1b[44mim\x1b[49m",
			"\x1b[44mbricated styles\x1b[49m\x1b[31m and \x1b[39m",
			"\x1b[31mon                  \x1b[39m",
			"\x1b[31mmultiple\x1b[39m lines.",
		},
		formatText("awesome text <fg=red>with <bg=blue>imbricated styles</> and on\nmultiple</> lines.", width, ""),
	)

	// Test edge-cases
	assert.Equal([]string{""}, formatText("", width, ""))
	assert.Equal([]string{""}, formatText("<fg=red></>", width, ""))
	assert.Equal([]string{"qsdf"}, formatText("<fg=wrong>qsdf</>", width, ""))
	assert.Equal([]string{"<toto=titi>qsdf</fg=", "blue>"}, formatText("<toto=titi>qsdf</fg=blue>", width, ""))
	assert.Equal([]string{"\x1b[34mtesttest\x1b[39m"}, formatText("<fg=blue>testtest", width, ""))
	assert.Equal([]string{"testt</fg=blue>est"}, formatText("testt</fg=blue>est", width, ""))
}

// TestFormatTextWithDefault checks that we can format text with a default style
func TestFormatTextWithDefault(t *testing.T) {
	assert := assert.New(t)
	width := 20

	assert.Equal(
		[]string{"\x1b[31;42mawesome text\x1b[39;49m"},
		formatText("<fg=red>awesome text</>", width, "bg=green;fg=blue"),
	)
	assert.Equal(
		[]string{"\x1b[34;42mawesome \x1b[39;49m\x1b[31;42mtext\x1b[39;49m"},
		formatText("awesome <fg=red>text</>", width, "bg=green;fg=blue"),
	)
	assert.Equal(
		[]string{
			"\x1b[34;42mawesome \x1b[39;49m\x1b[31;42mtext\x1b[39;49m\x1b[34;42m        \x1b[39;49m",
			"\x1b[34;42mwith \x1b[39;49m\x1b[33;42mmultiple lines\x1b[39;49m",
		},
		formatText("awesome <fg=red>text</>\nwith <fg=yellow>multiple lines</>", width, "bg=green;fg=blue"),
	)
}
