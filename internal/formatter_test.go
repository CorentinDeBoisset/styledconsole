package internal

import (
	"testing"

	"github.com/corentindeboisset/styledconsole/internal/style"
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

	assert.Equal("noop", EscapeTrailingBackslash("noop"))
	assert.Equal("super super\x00\x00", EscapeTrailingBackslash("super super\\\\"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00", EscapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00"))
	assert.Equal("super \x00 soupaire \x00 awesome \x00\x00\x00", EscapeTrailingBackslash("super \x00 soupaire \x00 awesome \x00\\\\"))
}

// TestFormatWithStyleWithoutTextBefore checks there are no errors of newline
func TestFormatWithStyleWithoutTextBefore(t *testing.T) {
	assert := assert.New(t)
	width := 20
	stack := style.OutputStyleStack{}

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
	stack.Push(*style.NewOutputStyle("bg=green;fg=red"))
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
		FormatText("great text", width, nil),
	)
	assert.Equal(
		[]string{"awesome text        ", "on                  ", "multiple lines."},
		FormatText("awesome text\non\nmultiple lines.", width, nil),
	)
	assert.Equal(
		[]string{"\x1b[31mawesome text\x1b[39m"},
		FormatText("<fg=red>awesome text</>", width, nil),
	)
	assert.Equal(
		[]string{
			"awesome text \x1b[31mwith st\x1b[39m",
			"\x1b[31myle and on          \x1b[39m",
			"\x1b[31mmultiple\x1b[39m lines.",
		},
		FormatText("awesome text <fg=red>with style and on\nmultiple</> lines.", width, nil),
	)
	assert.Equal(
		[]string{
			"awesome text \x1b[31mwith \x1b[39m\x1b[44mim\x1b[49m",
			"\x1b[44mbricated styles\x1b[49m\x1b[31m and \x1b[39m",
			"\x1b[31mon                  \x1b[39m",
			"\x1b[31mmultiple\x1b[39m lines.",
		},
		FormatText("awesome text <fg=red>with <bg=blue>imbricated styles</> and on\nmultiple</> lines.", width, nil),
	)

	// Test edge-cases
	assert.Equal([]string{""}, FormatText("", width, nil))
	assert.Equal([]string{""}, FormatText("<fg=red></>", width, nil))
	assert.Equal([]string{"qsdf"}, FormatText("<fg=wrong>qsdf</>", width, nil))
	assert.Equal([]string{"<toto=titi>qsdf"}, FormatText("<toto=titi>qsdf</fg=blue>", width, nil))
	assert.Equal([]string{"\x1b[34mtesttest\x1b[39m"}, FormatText("<fg=blue>testtest", width, nil))
	assert.Equal([]string{"testtest"}, FormatText("testt</fg=blue>est", width, nil))
}

// TestFormatTextWithDefault checks that we can format text with a default style
func TestFormatTextWithDefault(t *testing.T) {
	assert := assert.New(t)
	width := 20
	baseStyle := style.NewOutputStyle("bg=green;fg=blue")

	assert.Equal(
		[]string{"\x1b[31;42mawesome text\x1b[39;49m"},
		FormatText("<fg=red>awesome text</>", width, baseStyle),
	)
	assert.Equal(
		[]string{"\x1b[34;42mawesome \x1b[39;49m\x1b[31;42mtext\x1b[39;49m"},
		FormatText("awesome <fg=red>text</>", width, baseStyle),
	)
	assert.Equal(
		[]string{
			"\x1b[34;42mawesome \x1b[39;49m\x1b[31;42mtext\x1b[39;49m\x1b[34;42m                    \x1b[39;49m",
			"\x1b[34;42mwith \x1b[39;49m\x1b[33;42mmultiple lines\x1b[39;49m",
		},
		FormatText("awesome <fg=red>text</>\nwith <fg=yellow>multiple lines</>", width, baseStyle),
	)
}
