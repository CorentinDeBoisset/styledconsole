package style

// OutputStyleStack is a wrapper around a list of styles to have a first-in-last-out list behavior
type OutputStyleStack struct {
	styles []OutputStyle
}

// Push adds a new style to the stack
func (s *OutputStyleStack) Push (newStyle OutputStyle) {
	s.styles = append(s.styles, newStyle)
}

// Pop removes all styles after the given one. If there are no match in the stack, nothing is done
func (s *OutputStyleStack) Pop (oldStyle OutputStyle) {
	for i := len(s.styles) - 1; i > 0; i-- {
		if s.styles[i].Apply(``) == oldStyle.Apply(``) {
			s.styles = s.styles[:i]
			return
		}
	}
}

// PopCurrent removes the latest style in the stack
func (s *OutputStyleStack) PopCurrent () {
	if len(s.styles) > 0 {
		s.styles = s.styles[:len(s.styles)-1]
	}
}

// GetCurrent returns the latest style in the stack
func (s *OutputStyleStack) GetCurrent() OutputStyle {
	if len(s.styles) == 0 {
		return OutputStyle{}
	}

	return s.styles[len(s.styles) - 1]
}
