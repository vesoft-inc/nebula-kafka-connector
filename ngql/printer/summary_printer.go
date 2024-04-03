/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.qls
 */

package printer

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/vesoft-inc/nebula-ng-tools/golang/pkg/internel/generated_code/v5.0.0/proto/graph"
)

type SummaryPrinter struct {
	writer   table.Writer
	fd       *os.File
	filename string
	WidthMax int
}

func NewSummaryPrinter() SummaryPrinter {
	writer := table.NewWriter()
	configTableWriter(&writer, false)
	return SummaryPrinter{
		writer: writer,
	}
}

func (p *SummaryPrinter) closeFile() {
	if p.fd != nil {
		if err := p.fd.Close(); err != nil {
			fmt.Printf("Close file %s failed, %s", p.filename, err.Error())
		}
		p.fd = nil
		p.filename = ""
	}
}

func (p *SummaryPrinter) Print(summary *graph.Summary) {
	s := p.renderPlanInfo(summary.PlanInfo, string(summary.Preamble))
	fmt.Println(s)
	fmt.Printf("Execution Plan (build time %d us, optimize time %d us), [Px] means pipeline-x and [S] means storage side",
		summary.BuildTimeUs,
		summary.OptimizeTimeUs)

	if p.fd != nil {
		fmt.Fprintln(p.fd, s)
	}

	p.closeFile()
}

func (p *SummaryPrinter) renderPlanInfo(plan *graph.PlanInfo, preamble string) string {
	p.writer.ResetHeaders()
	p.writer.ResetRows()
	columnToFieldMap := map[string]string{
		"id":          "Id",
		"name":        "Name",
		"details":     "Details",
		"time(ms)":    "TimeMs",
		"rows":        "Rows",
		"memory(KiB)": "MemoryKib",
		"blocked(ms)": "BlockedMs",
		"queued(ms)":  "QueuedMs",
		"consume(ms)": "ConsumeMs",
		"produce(ms)": "ProduceMs",
		"finish(ms)":  "FinishMs",
		"batches":     "Batches",
		"concurrency": "Concurrency",
		"others":      "OtherStatsJson",
	}
	columns := []string{}
	switch preamble {
	case "explain":
		columns = append(columns, "details")
	case "profile":
		columns = append(columns, "details", "time(ms)", "rows", "memory(KiB)", "blocked(ms)")
	case "profile_verbose":
		columns = append(columns, "details", "time(ms)", "rows", "memory(KiB)", "blocked(ms)",
			"queued(ms)", "consume(ms)", "produce(ms)", "finish(ms)", "batches", "concurrency", "others")
	}

	fields := []string{}
	for _, column := range columns {
		fields = append(fields, columnToFieldMap[column])
	}

	// The first column is always the "plan" column which displays the plan tree.
	header := table.Row{"plan"}
	for _, column := range columns {
		header = append(header, column)
	}
	p.writer.AppendHeader(table.Row(header))

	rows := []table.Row{}
	FormatPretty(plan, "", "", fields, &rows)

	var columnConfigs []table.ColumnConfig
	columnConfigs = append(columnConfigs, table.ColumnConfig{
		Name:  "plan",
		Align: text.AlignDefault,
	})
	for _, column := range columns {
		align := text.AlignLeft
		if isNumeric(plan, columnToFieldMap[column]) {
			align = text.AlignRight
		}
		columnConfigs = append(columnConfigs, table.ColumnConfig{
			Name:     column,
			Align:    align,
			WidthMax: p.WidthMax,
			Transformer: func(val interface{}) string {
				if v, ok := val.(string); ok {
					if p.WidthMax != 0 && len(v) > p.WidthMax {
						return v[:p.WidthMax-3] + "..."
					}
				}
				return fmt.Sprint(val)
			},
		})
	}
	p.writer.SetColumnConfigs(columnConfigs)
	p.writer.Style().Options.DrawBorder = true
	p.writer.AppendRows(rows)

	s := p.writer.Render()
	return s
}

const (
	dash       string = `├─`
	spacer     string = `│ `
	dashLast   string = `└─`
	spacerLast string = `  `
)

func FormatPretty(p *graph.PlanInfo, prefix, childPrefix string, fields []string, rows *[]table.Row) {
	row := table.Row{}
	row = append(row, prefix+string(p.Name))
	for _, field := range fields {
		row = append(row, formattedFieldValue(p, field))
	}
	*rows = append(*rows, row)
	if len(p.Children) > 0 {
		for i, child := range p.Children {
			if i == len(p.Children)-1 {
				FormatPretty(child, childPrefix+dashLast, childPrefix+spacerLast, fields, rows)
			} else {
				FormatPretty(child, childPrefix+dash, childPrefix+spacer, fields, rows)
			}
		}
	}
}

func isNumeric(p *graph.PlanInfo, fieldName string) bool {
	value := fieldValue(p, fieldName)
	switch value.(type) {
	case int64, float64:
		return true
	}
	return false
}

func formattedFieldValue(p *graph.PlanInfo, fieldName string) string {
	value := fieldValue(p, fieldName)
	if value != nil {
		return formattedValue(value)
	}
	return ""
}

func fieldValue(p *graph.PlanInfo, fieldName string) interface{} {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	value := v.FieldByName(fieldName)
	if value.IsValid() {
		return value.Interface()
	}
	panic(fmt.Sprintf("field %s not found in the struct PlanInfo", fieldName))
}

func formattedValue(value interface{}) string {
	switch v := value.(type) {
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		if v == 0 {
			return "0"
		}
		return fmt.Sprintf("%.3f", v)
	case []byte:
		if json.Valid(v) {
			if string(v) == "{}" {
				return ""
			}
			return string(v)
		}
		return string(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
