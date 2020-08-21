package style

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEmptyStack checks the behavior of the stack when we try to get the current style while it's empty
func TestEmptyStack(t *testing.T) {
	assert := assert.New(t)
	stack := OutputStyleStack{}

	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestAddStyles checks we can add styles to the stack, fetch it and pop it
func TestAddStyles(t *testing.T) {
	assert := assert.New(t)
	stack := OutputStyleStack{}

	stack.Push(OutputStyle{Foreground: "green"})
	stack.Push(OutputStyle{Foreground: "blue"})

	assert.Equal(OutputStyle{Foreground: "blue"}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{Foreground: "green"}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestPopStyles checks we can add a style to the stack, fetch it and pop it
func TestPopStyles(t *testing.T) {
	assert := assert.New(t)
	stack := OutputStyleStack{}

	stack.Push(OutputStyle{Foreground: "green"})
	stack.Push(OutputStyle{Foreground: "blue"})
	stack.Push(OutputStyle{Foreground: "green"})
	stack.Push(OutputStyle{Foreground: "red"})

	assert.Equal(OutputStyle{Foreground: "red"}, stack.GetCurrent())

	stack.Pop(OutputStyle{Foreground: "green"}) // removes the green style and above
	assert.Equal(OutputStyle{Foreground: "blue"}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{Foreground: "green"}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{}, stack.GetCurrent())
}

// TestBaseStyle checks we can define a base style that will be merged with the stack styles
func TestBaseStyle(t *testing.T) {
	assert := assert.New(t)
	baseStyle := OutputStyle{Background: "green", Foreground: "blue"}
	stack := OutputStyleStack{BaseStyle: &baseStyle}
	stack.Push(OutputStyle{Foreground: "red"})

	assert.Equal(OutputStyle{Background: "green", Foreground: "red"}, stack.GetCurrent())
	stack.PopCurrent()
	assert.Equal(OutputStyle{Background: "green", Foreground: "blue"}, stack.GetCurrent())
}
