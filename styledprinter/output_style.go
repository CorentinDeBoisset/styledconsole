package styledprinter

import (
	"fmt"
	"os"
	"regexp"
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
	styleRegexp     = regexp.MustCompile(`([^=]+)=([^;]+)(;|$)`)
	separatorRegexp = regexp.MustCompile(`([^,;]+)`)
)

// OutputStyle contains the required data to print special text on the console
// They reference the available styles juste above
type OutputStyle struct {
	foreground           string
	background           string
	href                 string
	handleHrefGracefully bool
	options              map[string]bool
}

func NewOutputStyle(styleString string) *OutputStyle {
	styleMatches := styleRegexp.FindAllSubmatchIndex([]byte(styleString), -1)
	if len(styleMatches) == 0 {
		return nil
	}

	var foreground string
	var background string
	var href string
	var handleHrefGracefully bool
	var options map[string]bool

	for _, attrMatches := range styleMatches {
		styleName := strings.ToLower(styleString[attrMatches[2]:attrMatches[3]])
		styleValue := styleString[attrMatches[4]:attrMatches[5]]

		if `fg` == styleName {
			foreground = strings.ToLower(styleValue)
		} else if `bg` == styleName {
			background = strings.ToLower(styleValue)
		} else if `href` == styleName {
			href = styleValue
		} else if `options` == styleName {
			if options == nil {
				options = make(map[string]bool)
			}
			separatorMatches := separatorRegexp.FindAllSubmatchIndex([]byte(strings.ToLower(styleValue)), -1)
			for _, separatorIndexes := range separatorMatches {
				options[styleValue[separatorIndexes[2]:separatorIndexes[3]]] = true
			}
		} else {
			// If there is an unknown attribute, the whole style is voided
			return nil
		}
	}

	handleHrefGracefully = os.Getenv("TERMINAL_EMULATOR") != `JetBrains-JediTerm` && os.Getenv("KONSOLE_VERSION") == ``

	return &OutputStyle{
		foreground:           foreground,
		background:           background,
		href:                 href,
		handleHrefGracefully: handleHrefGracefully,
		options:              options,
	}
}

// Apply surrounds a given string with the adequate escape sequence
func (s OutputStyle) Apply(text string) string {
	var setCodes []string
	var unsetCodes []string

	if foreground, ok := availableForegroundColors[s.foreground]; ok {
		setCodes = append(setCodes, foreground[0])
		unsetCodes = append(unsetCodes, foreground[1])
	}
	if background, ok := availableBackgroundColors[s.background]; ok {
		setCodes = append(setCodes, background[0])
		unsetCodes = append(unsetCodes, background[1])
	}
	for styleOption, enabled := range s.options {
		if option, ok := availableOptions[styleOption]; ok && enabled {
			setCodes = append(setCodes, option[0])
			unsetCodes = append(unsetCodes, option[1])
		}
	}

	if s.href != `` && s.handleHrefGracefully {
		return fmt.Sprintf("\033]8;;%s\033\\%s\033]8;;\033\\", s.href, text)
	}

	if len(setCodes) == 0 {
		return text
	}

	// Sort codes in order to always have the same output for two similar styles
	sort.Strings(setCodes)
	sort.Strings(unsetCodes)

	return fmt.Sprintf("\033[%sm%s\033[%sm", strings.Join(setCodes, ";"), text, strings.Join(unsetCodes, ";"))
}
