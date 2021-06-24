package styledconsole

import (
	"fmt"
	"strings"
)

func getColumnWidths(headers []string, content [][]string) []int {
	columnCount := len(headers)
	for _, row := range content {
		if columnCount < len(row) {
			columnCount = len(row)
		}
	}

	var columnWidths = make([]int, columnCount)

	for i, headerItem := range headers {
		itemWidth := 0
		for _, line := range strings.Split(headerItem, "\n") {
			if itemWidth < len(line) {
				itemWidth = len(line)
			}
		}
		columnWidths[i] = itemWidth
	}
	for _, row := range content {
		for i, rowItem := range row {
			itemWidth := 0
			for _, line := range strings.Split(rowItem, "\n") {
				if itemWidth < len(line) {
					itemWidth = len(line)
				}
			}
			if columnWidths[i] < itemWidth {
				columnWidths[i] = itemWidth
			}
		}
	}

	return columnWidths
}

func getTableWidth(widths []int) int {
	totalWidth := 0
	for _, width := range widths {
		totalWidth += width
	}
	totalWidth += 4 + 3*(len(widths)-1)

	return totalWidth
}

func getAcceptableColumnWidths(startingWidths []int, termWidth int) []int {
	// Deep copy of startingWidths
	columnWidths := make([]int, len(startingWidths))
	copy(columnWidths, startingWidths)

	for {
		if getTableWidth(columnWidths) < termWidth {
			// It fits, we can stop
			return columnWidths
		}

		// Every time, we try to reduce the largest column
		largestIdx := 0
		largestWidth := 0
		for i, wid := range columnWidths {
			if largestWidth < wid {
				largestWidth = wid
				largestIdx = i
			}
		}
		if largestWidth > 15 {
			if getTableWidth(columnWidths)-termWidth < largestWidth-15 {
				// If we reduce the largest column, it will fit
				columnWidths[largestIdx] -= getTableWidth(columnWidths) - termWidth
				return columnWidths
			} else {
				// We reduce the largest column to 15char, it won't fit so we continue
				columnWidths[largestIdx] = 15
				continue
			}
		}

		// If we cannot reduce any column, we try some last resort option
		avgWidth := int(termWidth / len(columnWidths))
		if avgWidth >= 10 {
			for i := range columnWidths {
				columnWidths[i] = avgWidth
			}
		}

		// The total width is too wide, but we did our best...
		return columnWidths
	}
}

func formatOneRow(row []string, columnWidths []int) string {
	var preparedSubLines [][]string
	totalLines := 1

	for cellIdx, cell := range row {
		var preparedCellLines []string
		for _, subLine := range strings.Split(cell, "\n") {
			if len(subLine) <= columnWidths[cellIdx] {
				preparedCellLines = append(preparedCellLines, subLine)
			} else {
				i := 0
				for {
					if i+columnWidths[cellIdx] > len(subLine) {
						preparedCellLines = append(preparedCellLines, subLine[i:])
						break
					} else {
						preparedCellLines = append(preparedCellLines, subLine[i:i+columnWidths[cellIdx]])
						i += columnWidths[cellIdx]
					}
				}
			}
		}
		if len(preparedCellLines) > totalLines {
			totalLines = len(preparedCellLines)
		}
		preparedSubLines = append(preparedSubLines, preparedCellLines)
	}

	var rowsToPrint []string
	for i := 0; i < totalLines; i++ {
		rowToPrint := "|"
		for columnIdx, column := range preparedSubLines {
			if i < len(column) {
				rowToPrint += fmt.Sprintf(" %s%s |", column[i], strings.Repeat(" ", columnWidths[columnIdx]-len(column[i])))
			} else {
				rowToPrint += fmt.Sprintf(" %s |", strings.Repeat(" ", columnWidths[columnIdx]))
			}
		}
		rowsToPrint = append(rowsToPrint, rowToPrint)
	}

	return strings.Join(rowsToPrint, "\n")
}
