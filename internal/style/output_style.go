package style

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

var (
	// All the available styles are in the form [set, unset]
	// where the int value are the escape sequence to send in the console
	// to either set or unset the style
	availableForegroundColors = map[string]([2]string){
		"black":   [2]string{"30", "39"},
		"red":     [2]string{"31", "39"},
		"green":   [2]string{"32", "39"},
		"yellow":  [2]string{"33", "39"},
		"blue":    [2]string{"34", "39"},
		"magenta": [2]string{"35", "39"},
		"cyan":    [2]string{"36", "39"},
		"white":   [2]string{"37", "39"},
		"default": [2]string{"39", "39"},
	}
	availableBackgroundColors = map[string]([2]string){
		"black":   [2]string{"40", "49"},
		"red":     [2]string{"41", "49"},
		"green":   [2]string{"42", "49"},
		"yellow":  [2]string{"43", "49"},
		"blue":    [2]string{"44", "49"},
		"magenta": [2]string{"45", "49"},
		"cyan":    [2]string{"46", "49"},
		"white":   [2]string{"47", "49"},
		"default": [2]string{"49", "49"},
	}
	availableOptions = map[string]([2]string){
		"bold":       [2]string{"1", "22"},
		"underscore": [2]string{"4", "24"},
		"blink":      [2]string{"5", "25"},
		"reverse":    [2]string{"7", "27"},
		"conceal":    [2]string{"8", "28"},
	}
)

// OutputStyle contains the required data to print special text on the console
// They reference the available styles juste above
type OutputStyle struct {
	Foreground           string
	Background           string
	Href                 string
	HandleHrefGracefully *bool
	Options              []string
}

// Apply surrounds a given string with the adequate escape sequence
func (s OutputStyle) Apply(text string) string {
	var setCodes []string
	var unsetCodes []string

	if s.HandleHrefGracefully == nil {
		s.HandleHrefGracefully = new(bool)
		*s.HandleHrefGracefully = os.Getenv("TERMINAL_EMULATOR") != `JetBrains-JediTerm` && os.Getenv("KONSOLE_VERSION") == ``
	}

	if foreground, ok := availableForegroundColors[s.Foreground]; ok {
		setCodes = append(setCodes, foreground[0])
		unsetCodes = append(unsetCodes, foreground[1])
	}
	if background, ok := availableBackgroundColors[s.Background]; ok {
		setCodes = append(setCodes, background[0])
		unsetCodes = append(unsetCodes, background[1])
	}
	for _, styleOption := range s.Options {
		if option, ok := availableOptions[styleOption]; ok {
			setCodes = append(setCodes, option[0])
			unsetCodes = append(unsetCodes, option[1])
		}
	}

	if s.Href != `` && s.HandleHrefGracefully != nil && *s.HandleHrefGracefully {
		return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", s.Href, text)
	}

	if len(setCodes) == 0 {
		return text
	}

	// Sort codes in order to always have the same output for two similar styles
	sort.Strings(setCodes)
	sort.Strings(unsetCodes)

	return fmt.Sprintf("\033[%sm%s\033[%sm", strings.Join(setCodes, ";"), text, strings.Join(unsetCodes, ";"))
}
