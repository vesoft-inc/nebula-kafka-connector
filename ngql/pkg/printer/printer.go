package printer

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang/pkg/types"
)

const (
	dash       string = `├─`
	spacer     string = `│ `
	dashLast   string = `└─`
	spacerLast string = `  `
)

type Printer interface {
	SetFormat(string)
	PrintResult(io.Writer, nebula.Result)
	// export the result to csv format
	ExportResultCsv(io.Writer, nebula.Result)
	PrintResultVertical(io.Writer, nebula.Result)
	PrintPlanInfo(w io.Writer, summary nebula.Summary)
}

type defaultPrinter struct {
	format     string
	tableWrite table.Writer
	stringer   valueStringer
	widthMax   int
}

func configTableWriter(writer table.Writer, separateRows bool) {
	writer.Style().Format.Header = text.FormatDefault
	writer.Style().Options.SeparateRows = separateRows
}

func NewPrinter(format string, widthMax int) Printer {
	writer := table.NewWriter()
	stringer := newValueStringer(format)
	configTableWriter(writer, false)
	return &defaultPrinter{
		format:     format,
		tableWrite: writer,
		stringer:   stringer,
		widthMax:   widthMax,
	}
}

func (p *defaultPrinter) SetFormat(format string) {
	p.format = format
	p.stringer = newValueStringer(format)
}

func (p *defaultPrinter) PrintResult(w io.Writer, res nebula.Result) {
	var s string
	s += p.getResultString(res)
	if p.getQueryStatsString(res) != "" {
		s += "\n"
		s += p.getQueryStatsString(res)
	}
	fmt.Fprintln(w, s)
}

func (p *defaultPrinter) ExportResultCsv(w io.Writer, res nebula.Result) {
	s := p.getResultCsvString(res)
	fmt.Fprintln(w, s)
}
func (p *defaultPrinter) PrintResultVertical(w io.Writer, res nebula.Result) {
	var s string
	s += p.getResultVerticalString(res)
	if p.getQueryStatsString(res) != "" {
		s += p.getQueryStatsString(res)
	}
	fmt.Fprintln(w, s)
}

func (p *defaultPrinter) PrintPlanInfo(w io.Writer, summary nebula.Summary) {
	s := p.renderPlanInfo(summary.PlanInfo(), string(summary.ExplainType()))
	fmt.Fprintln(w, s)
	fmt.Fprintf(w, "Execution Plan: [Px] means pipeline-x and [S] means storage side.\n\n")
	fmt.Fprintf(w, "Elapsed Time:\n")
	fmt.Fprintf(w, " build time    : %d us\n", summary.BuildTimeUs())
	fmt.Fprintf(w, " optimize time : %d us\n", summary.OptimizeTimeUs())
	fmt.Fprintf(w, " serialize time: %d us\n", summary.SerializeTimeUs())
	fmt.Fprintf(w, " total time    : %d us\n", summary.TotalServerTimeUs())
	fmt.Fprintln(w, "")
}

// use the console string() format
func (p *defaultPrinter) getResultString(res nebula.Result) string {
	var s string
	if len(res.Columns()) != 0 {
		p.tableWrite.ResetHeaders()
		p.tableWrite.ResetRows()
		header := make([]interface{}, 0)
		for _, col := range res.Columns() {
			header = append(header, col)
		}
		p.tableWrite.AppendHeader(header)
		for res.HasNext() {
			row, err := res.Next()
			if err != nil {
				return "INVALID ROW, err: " + err.Error()
			}
			var values []interface{}
			for _, col := range row.Values() {
				values = append(values, p.stringer.String(col))
			}
			p.tableWrite.AppendRow(values)
		}
		s += p.tableWrite.Render()
	}
	return s
}

// use the default string() in nebula-go
func (p *defaultPrinter) getResultCsvString(res nebula.Result) string {
	var s string

	if len(res.Columns()) != 0 {
		p.tableWrite.ResetHeaders()
		p.tableWrite.ResetRows()
		header := make([]interface{}, 0)
		for _, col := range res.Columns() {
			header = append(header, col)
		}
		p.tableWrite.AppendHeader(header)
		for res.HasNext() {
			row, err := res.Next()
			if err != nil {
				continue
			}
			var values []interface{}
			for _, val := range row.Values() {
				values = append(values, val.String())
			}
			p.tableWrite.AppendRow(values)
		}
		s += p.tableWrite.Render()
	}
	return s
}

func (p *defaultPrinter) getQueryStatsString(res nebula.Result) string {
	var s string
	stats := res.Summary().QueryStats()
	// TODO(jie): Add a QueryType field to Summary to specify if the query is a DML.
	// Only print the mentioned stats for DML queries.
	if stats.NumAffectedNodes() > 0 {
		s += fmt.Sprintf("Affected nodes: %d\n", stats.NumAffectedNodes())
	}
	if stats.NumAffectedEdges() > 0 {
		s += fmt.Sprintf("Affected edges: %d\n", stats.NumAffectedEdges())
	}
	return s
}

func (p *defaultPrinter) getResultVerticalString(res nebula.Result) string {
	var s string
	if res.RowSize() == 0 {
		return ""
	}

	maxColumnWidth := 0
	for _, columName := range res.Columns() {
		if len(columName) > maxColumnWidth {
			maxColumnWidth = len(columName)
		}
	}
	numCols := len(res.Columns())
	var i int
	for res.HasNext() {
		s += fmt.Sprintf("*************************** %d. row ***************************\n", i+1)
		record, err := res.Next()
		if err != nil {
			continue
		}
		i++
		for j := 0; j < numCols; j++ {
			val, err := record.GetValueByIndex(j)
			if err != nil {
				continue
			}
			colName := res.Columns()[j]
			s += fmt.Sprintf("%*s: %s\n", maxColumnWidth, colName, p.stringer.String(val))
		}
	}
	return s

}

func (p *defaultPrinter) renderPlanInfo(plan nebula.PlanInfo, explainType string) string {
	p.tableWrite.ResetHeaders()
	p.tableWrite.ResetRows()
	columnToFieldMap := map[string]string{
		"id":          "Id",
		"name":        "Name",
		"details":     "Details",
		"columns":     "Columns",
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
	isVerbose := false
	columns := []string{}
	switch explainType {
	case "explain":
		columns = append(columns, "details")
	case "explain_verbose":
		isVerbose = true
		columns = append(columns, "details", "columns")
	case "profile":
		columns = append(columns, "details", "time(ms)", "rows", "memory(KiB)", "blocked(ms)")
	case "profile_verbose":
		isVerbose = true
		columns = append(columns, "details", "time(ms)", "rows", "memory(KiB)", "blocked(ms)",
			"queued(ms)", "consume(ms)", "produce(ms)", "finish(ms)", "batches", "concurrency", "others")
	default:
		if explainType != "" {
			panic(fmt.Sprintf("unknown explainType: %s", explainType))
		}
	}

	widthMax := p.widthMax
	if isVerbose {
		// If the verbose mode is enabled, we will not truncate the columns when printing the plan info.
		widthMax = 0
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
	p.tableWrite.AppendHeader(table.Row(header))

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
			WidthMax: widthMax,
			Transformer: func(val interface{}) string {
				if v, ok := val.(string); ok {
					if widthMax != 0 && len(v) > p.widthMax {
						return v[:widthMax-3] + "..."
					}
				}
				return fmt.Sprint(val)
			},
		})
	}
	p.tableWrite.SetColumnConfigs(columnConfigs)
	p.tableWrite.Style().Options.DrawBorder = true
	p.tableWrite.AppendRows(rows)

	s := p.tableWrite.Render()
	return s
}

func FormatPretty(p nebula.PlanInfo, prefix, childPrefix string, fields []string, rows *[]table.Row) {
	row := table.Row{}
	row = append(row, prefix+string(p.Name()))
	for _, field := range fields {
		row = append(row, formattedFieldValue(p, field))
	}
	*rows = append(*rows, row)
	if len(p.Children()) > 0 {
		for i, child := range p.Children() {
			if i == len(p.Children())-1 {
				FormatPretty(child, childPrefix+dashLast, childPrefix+spacerLast, fields, rows)
			} else {
				FormatPretty(child, childPrefix+dash, childPrefix+spacer, fields, rows)
			}
		}
	}
}

func isNumeric(p nebula.PlanInfo, fieldName string) bool {
	value := fieldValue(p, fieldName)
	switch value.(type) {
	case int64, float64:
		return true
	}
	return false
}

func formattedFieldValue(p nebula.PlanInfo, fieldName string) string {
	value := fieldValue(p, fieldName)
	if value != nil {
		return formattedValue(value)
	}
	return ""
}

func fieldValue(p nebula.PlanInfo, fieldName string) interface{} {
	v := reflect.ValueOf(p)

	method := v.MethodByName(fieldName)
	if !method.IsValid() {
		panic(fmt.Sprintf("method %s not found in the interface PlanInfo", fieldName))
	}
	value := method.Call([]reflect.Value{})
	if len(value) == 0 {
		return nil
	}
	// only return the first value if it is valid
	if value[0].IsValid() {
		return value[0].Interface()
	}
	panic(fmt.Sprintf("field %s not found in the interface PlanInfo", fieldName))
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
	case []string:
		return fmt.Sprintf("[%s]", strings.Join(v, ", "))
	default:
		return fmt.Sprintf("%v", v)
	}
}
