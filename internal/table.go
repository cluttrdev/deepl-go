package table

import (
	"fmt"

	"github.com/pkg/errors"
)

type Table struct {
	Headers []string
	Rows    [][]string
}

func NewTable(headers ...string) (t *Table) {
	return &Table{
		Headers: headers,
		Rows:    make([][]string, 0),
	}
}

func (t *Table) AddRow(cols ...string) error {
	if len(cols) != len(t.Headers) {
		return errors.Errorf("Number of columns (%d) does not match headers (%d)", len(cols), len(t.Headers))
	}

	t.Rows = append(t.Rows, cols)
	return nil
}

func (t *Table) Print() {
	colWidths := make([]int, len(t.Headers))

	for j, h := range t.Headers {
		colWidths[j] = len(h)
	}

	for _, r := range t.Rows {
		for j, c := range r {
			if len(c) > colWidths[j] {
				colWidths[j] = len(c)
			}
		}
	}

	for j, h := range t.Headers {
		fmt.Printf("%-*s  ", colWidths[j], h)
	}
	fmt.Println()
	for _, r := range t.Rows {
		for j, c := range r {
			fmt.Printf("%-*s  ", colWidths[j], c)
		}
		fmt.Println()
	}
}
