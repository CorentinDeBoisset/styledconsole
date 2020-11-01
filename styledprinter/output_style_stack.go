package styledprinter

// outputStyleStack is a wrapper around a list of styles to have a first-in-last-out list behavior
type outputStyleStack struct {
	baseStyleString string
	baseStyle       *OutputStyle
	styles          []OutputStyle
}

func newOutputStyleStack(baseStyleString string) outputStyleStack {
	if baseStyleString != "" {
		baseStyle := NewOutputStyle(baseStyleString)
		if baseStyle != nil {
			return outputStyleStack{
				baseStyleString: baseStyleString,
				baseStyle:       baseStyle,
			}
		}
	}

	return outputStyleStack{}
}

// Push adds a new style to the stack
func (s *outputStyleStack) Push(newStyleString string) bool {
	var newStyle *OutputStyle
	if s.baseStyleString != "" {
		newStyle = NewOutputStyle(s.baseStyleString + ";" + newStyleString)
	} else {
		newStyle = NewOutputStyle(newStyleString)
	}

	if newStyle != nil {
		s.styles = append(s.styles, *newStyle)
		return true
	}

	return false
}

// Pop removes all styles after the given one. If there are no match in the stack, nothing is done
func (s *outputStyleStack) Pop(oldStyleString string) bool {
	oldStyle := NewOutputStyle(oldStyleString)
	if oldStyle == nil {
		return false
	}

	for i := len(s.styles) - 1; i > 0; i-- {
		if s.styles[i].Apply(``) == oldStyle.Apply(``) {
			s.styles = s.styles[:i]
			return true
		}
	}

	return false
}

// PopCurrent removes the latest style in the stack
func (s *outputStyleStack) PopCurrent() {
	if len(s.styles) > 0 {
		s.styles = s.styles[:len(s.styles)-1]
	}
}

// GetCurrent returns the latest style in the stack
func (s *outputStyleStack) GetCurrent() OutputStyle {
	if len(s.styles) == 0 {
		if s.baseStyle != nil {
			return *s.baseStyle
		}

		return OutputStyle{}
	}

	return s.styles[len(s.styles)-1]
}
