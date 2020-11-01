package styledconsole

import (
	"github.com/corentindeboisset/styledconsole/styledprinter"
)

var greenStyle, yellowStyle, redStyle, highlightedChoiceStyle *styledprinter.OutputStyle

func init() {
	greenStyle = styledprinter.NewOutputStyle("fg=green")
	redStyle = styledprinter.NewOutputStyle("fg=red")
	highlightedChoiceStyle = styledprinter.NewOutputStyle("fg=cyan;options=bold,underscore")
	yellowStyle = styledprinter.NewOutputStyle("fg=yellow")
}
