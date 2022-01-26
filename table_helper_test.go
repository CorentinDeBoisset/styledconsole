package styledconsole

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetSubstring checks that column widths are correct
func TestGetColumnWidths(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(
		[]int{7, 7, 5},
		getColumnWidths(
			[]string{"one", "two", "three"},
			[][]string{{"1", "big two", "3"}, {"big one", "2", "3"}},
		),
	)

	// with UTF8
	assert.Equal(
		[]int{7, 7, 5},
		getColumnWidths(
			[]string{"un", "deux", "troïs"},
			[][]string{{"1", "big twô", "3"}, {"big @ne", "2", "3"}},
		),
	)
}

// TestFormatOneRow checks rows are well formatted
func TestFormatOneRow(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(
		"| one     | two     | three |",
		formatOneRow(
			[]string{"one", "two", "three"},
			[]int{7, 7, 5},
		),
	)

	// with UTF8
	assert.Equal(
		"| un      | deux    | troïs |",
		formatOneRow(
			[]string{"un", "deux", "troïs"},
			[]int{7, 7, 5},
		),
	)
}
