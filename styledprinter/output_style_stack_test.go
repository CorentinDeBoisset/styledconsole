package styledprinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEmptyStack checks the behavior of the stack when we try to get the current style while it's empty
func TestEmptyStack(t *testing.T) {
	assert := assert.New(t)
	stack := newOutputStyleStack("")

	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestAddStyles checks we can add styles to the stack, fetch it and pop it
func TestAddStyles(t *testing.T) {
	assert := assert.New(t)
	stack := newOutputStyleStack("")

	stack.Push("fg=green")
	stack.Push("fg=blue")

	assert.Equal(OutputStyle{foreground: "blue", handleHrefGracefully: true}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{foreground: "green", handleHrefGracefully: true}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestPopStyles checks we can add a style to the stack, fetch it and pop it
func TestPopStyles(t *testing.T) {
	assert := assert.New(t)
	stack := newOutputStyleStack("")

	stack.Push("fg=green")
	stack.Push("fg=blue")
	stack.Push("fg=green")
	stack.Push("fg=red")

	assert.Equal(OutputStyle{foreground: "red", handleHrefGracefully: true}, stack.GetCurrent())

	stack.Pop("fg=green") // removes the green style and above
	assert.Equal(OutputStyle{foreground: "blue", handleHrefGracefully: true}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{foreground: "green", handleHrefGracefully: true}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestBaseStyle checks we can define a base style that will be merged with the stack styles
func TestBaseStyle(t *testing.T) {
	assert := assert.New(t)
	stack := newOutputStyleStack("bg=green;fg=blue")
	stack.Push("fg=red")

	assert.Equal(OutputStyle{background: "green", foreground: "red", handleHrefGracefully: true}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{background: "green", foreground: "blue", handleHrefGracefully: true}, stack.GetCurrent())
}
