package printer

import (
	"bytes"
	"fmt"
)

func FormatTable(headers []string, data [][]string) string {
	var buf bytes.Buffer

	// print header
	for _, header := range headers {
		fmt.Fprintf(&buf, "%-20s", header)
	}
	fmt.Fprintln(&buf)

	// print split line
	for i := 0; i < len(headers)*20; i++ {
		fmt.Fprint(&buf, "-")
	}
	fmt.Fprintln(&buf)

	// print data
	for _, row := range data {
		for _, cell := range row {
			fmt.Fprintf(&buf, "%-20s", cell)
		}
		fmt.Fprintln(&buf)
	}

	return buf.String()
}
