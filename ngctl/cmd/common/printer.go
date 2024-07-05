package common

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func FormatTable(headers []string, data [][]string) string {
	writer := table.NewWriter()
	writer.Style().Format.Header = text.FormatTitle
	writer.Style().Options.DrawBorder = false
	writer.Style().Options.SeparateColumns = false
	writer.Style().Options.SeparateHeader = false

	writer.ResetHeaders()
	writer.ResetRows()
	header := make([]interface{}, 0)
	for _, h := range headers {
		header = append(header, h)
	}
	writer.AppendHeader(header)
	for _, d := range data {
		var values []interface{}
		for _, val := range d {
			values = append(values, val)
		}
		writer.AppendRow(values)
	}
	return writer.Render()
}
