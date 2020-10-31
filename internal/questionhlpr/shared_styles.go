package questionhlpr

import "github.com/corentindeboisset/styledconsole/internal/style"

var greenStyle, yellowStyle, redStyle, highlightedChoiceStyle *style.OutputStyle

func init() {
	greenStyle = style.NewOutputStyle("fg=green")
	redStyle = style.NewOutputStyle("fg=red")
	highlightedChoiceStyle = style.NewOutputStyle("fg=cyan;options=bold,underscore")
	yellowStyle = style.NewOutputStyle("fg=yellow")
}
