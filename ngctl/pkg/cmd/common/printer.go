package common

import (
	"encoding/json"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type OutputFormatType string

const (
	FormatTypeTable OutputFormatType = "table"
	FormatTypeJson                   = "json"
)

func Format(headers []string, data [][]string, typ OutputFormatType) (string, error) {
	switch typ {
	case FormatTypeTable:
		return formatTable(headers, data)
	case FormatTypeJson:
		return formatJson(headers, data)
	default:
		return "", fmt.Errorf("unsupported format type, %s", typ)
	}
}

func formatJson(headers []string, data [][]string) (string, error) {
	m := make([]map[string]string, 0)
	for _, d := range data {
		row := make(map[string]string)
		for i, val := range d {
			row[headers[i]] = val
		}
		m = append(m, row)
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}

func formatTable(headers []string, data [][]string) (string, error) {
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
	return writer.Render(), nil
}
