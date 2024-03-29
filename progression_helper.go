package styledconsole

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

var progressBarLength int = 19
var progressStarted bool = false
var progressTotalSteps int
var progressDone int
var lastPrintAdvancement int
var lastPrintTime time.Time

// ProgressStart starts a progress bar of a given duration
func ProgressStart(totalSteps int) {
	if progressStarted || totalSteps < 0 {
		return
	}

	fmt.Printf("  %s", buildProgressBar(0, totalSteps))
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Print("\n")
	}

	progressStarted = true
	progressTotalSteps = totalSteps
	progressDone = 0
	lastPrintAdvancement = 0
	lastPrintTime = time.Now()
}

// ProgressAdvance advances the current progress bar of a given stepCount. If there is no progressBar it does nothing
func ProgressAdvance(stepCount int) {
	if !progressStarted || stepCount == 0 {
		return
	}

	if progressDone+stepCount < progressTotalSteps {
		progressDone += stepCount
		if term.IsTerminal(int(os.Stdout.Fd())) {
			if time.Since(lastPrintTime) > 1e8 {
				// Do not refresh faster than 10fps, to avoid stdout-induced lag
				fmt.Printf("\033[1000D  %s", buildProgressBar(progressDone, progressTotalSteps))
				lastPrintTime = time.Now()
			}
		} else if (progressDone - lastPrintAdvancement) >= int(float64(progressTotalSteps)*0.05) {
			fmt.Printf("  %s\n", buildProgressBar(progressDone, progressTotalSteps))
			lastPrintAdvancement = progressDone
		}
	} else {
		progressDone = progressTotalSteps
		ProgressFinish()
	}
}

// ProgressFinish finishes the current progress bar. If there is no progressBar it does nothing
func ProgressFinish() {
	if !progressStarted {
		return
	}

	if term.IsTerminal(int(os.Stdout.Fd())) {
		fmt.Printf("\033[1000D  %s\n", buildProgressBar(progressTotalSteps, progressTotalSteps))
	} else {
		fmt.Printf("  %s\n", buildProgressBar(progressTotalSteps, progressTotalSteps))
	}

	progressStarted = false
}

func buildProgressBar(done int, total int) string {
	advancementRatio := float64(done) / float64(total)
	advancement := int(math.Round(advancementRatio * float64(progressBarLength)))

	maxDigits := int(math.Ceil(math.Log10(float64(total + 1))))
	paddedDone := fmt.Sprintf("%-"+strconv.Itoa(maxDigits)+"d", done)

	return fmt.Sprintf(
		"%s/%d [%s>%s] %3d%%",
		paddedDone,
		total,
		strings.Repeat("=", advancement),
		strings.Repeat("-", progressBarLength-advancement),
		int(advancementRatio*100),
	)
}
