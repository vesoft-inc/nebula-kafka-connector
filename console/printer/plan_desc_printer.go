/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.qls
 */

package printer

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

func graphvizString(s string) string {
	s = strings.Replace(s, "{", "\\{", -1)
	s = strings.Replace(s, "}", "\\}", -1)
	s = strings.Replace(s, "\"", "\\\"", -1)
	s = strings.Replace(s, "[", "\\[", -1)
	s = strings.Replace(s, "]", "\\]", -1)
	return s
}

type PlanDescPrinter struct {
	writer   table.Writer
	fd       *os.File
	filename string
	WidthMax int
}

func NewPlanDescPrinter() PlanDescPrinter {
	writer := table.NewWriter()
	configTableWriter(&writer, false)
	return PlanDescPrinter{
		writer: writer,
	}
}

func (p *PlanDescPrinter) ExportDot(filename string) {
	if filename == "" {
		p.closeFile()
		return
	}
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Open or Create file %s failed, %s", filename, err.Error())
		return
	}
	p.fd = fd
	p.filename = filename
}

func (p *PlanDescPrinter) closeFile() {
	if p.fd != nil {
		if err := p.fd.Close(); err != nil {
			fmt.Printf("Close file %s failed, %s", p.filename, err.Error())
		}
		p.fd = nil
		p.filename = ""
	}
}

func (p PlanDescPrinter) configWriterDotRenderStyle(renderByDot bool) {
	if renderByDot {
		p.writer.Style().Box.Left = " "
		p.writer.Style().Box.Right = " "
	} else {
		p.writer.Style().Box.Left = "|"
		p.writer.Style().Box.Right = "|"
	}
	p.writer.Style().Box.BottomLeft = "-"
	p.writer.Style().Box.BottomRight = "-"
	p.writer.Style().Box.TopLeft = "-"
	p.writer.Style().Box.TopRight = "-"
	p.writer.Style().Box.LeftSeparator = "-"
	p.writer.Style().Box.RightSeparator = "-"
}

func (p PlanDescPrinter) renderDotGraph(s string) string {
	p.writer.ResetHeaders()
	p.writer.ResetRows()
	p.configWriterDotRenderStyle(true)
	p.writer.AppendHeader(table.Row{"plan"})
	p.writer.AppendRow(table.Row{s})
	return p.writer.Render()
}

func (p PlanDescPrinter) renderDotGraphByStruct(s string) string {
	p.writer.ResetHeaders()
	p.writer.ResetRows()
	p.configWriterDotRenderStyle(true)
	p.writer.AppendHeader(table.Row{"plan"})
	p.writer.AppendRow(table.Row{s})
	return p.writer.Render()
}

func (p PlanDescPrinter) renderByRow(rs *nebula.ResultSet) string {
	p.writer.ResetHeaders()
	p.writer.ResetRows()
	p.configWriterDotRenderStyle(false)
	var columnConfigs []table.ColumnConfig
	headerRow := table.Row{}
	rightSepToTailWidth, rows := rs.MakePlanByRow()
	header := rs.GetHeader()

	if len(rightSepToTailWidth) != len(header) {
		log.Fatalf("rightSepToTailWidth and header length not equal: %d vs %d", len(rightSepToTailWidth), len(header))
	}

	for i, col := range header {
		headerRow = append(headerRow, col)
		if strings.Contains(col, "/") {
			width := rightSepToTailWidth[i]
			columnConfigs = append(columnConfigs, table.ColumnConfig{
				Name:  col,
				Align: text.AlignRight,
				Transformer: func(s interface{}) string {
					if v, ok := s.(string); ok && width > 0 && strings.Contains(v, "/") {
						l := len(v) - strings.LastIndex(v, "/")
						if width < l {
							return v
						}
						return fmt.Sprintf("%s%s", v, strings.Repeat(" ", width-l))
					}
					return fmt.Sprint(s)
				},
			})
		} else if i != 0 {
			align := text.AlignLeft
			if col == "time" || col == "memory" {
				align = text.AlignRight
			}
			columnConfigs = append(columnConfigs, table.ColumnConfig{
				Name:     col,
				Align:    align,
				WidthMax: p.WidthMax,
				// WidthMaxEnforcer: func(s string, wrapLen int) string {
				// 	if len(s) > wrapLen {
				// 		s = s[:wrapLen-3] + "..."
				// 	}
				// 	return s
				// },
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
	}
	p.writer.AppendHeader(headerRow)
	if len(columnConfigs) > 0 {
		p.writer.SetColumnConfigs(columnConfigs)
	}

	for _, row := range rows {
		p.writer.AppendRow(table.Row(row))
	}
	return p.writer.Render()
}

func (p *PlanDescPrinter) PrintPlanDesc(res *nebula.ResultSet) {
	var s string
	format := strings.ToLower(res.GetPlanPrintFormat())
	switch format {
	case "row":
		s = p.renderByRow(res)
		fmt.Println(s)
	case "dot":
		s = res.MakeDotGraph()
		fmt.Println(p.renderDotGraph(s))
	case "dot:struct":
		s = res.MakeDotGraphByStruct()
		fmt.Println(p.renderDotGraphByStruct(s))
	}
	fmt.Printf("Execution Plan (build time %d us, optimize time %d us), [Px] means pipeline-x and [S] means storage side",
		res.GetBuildTimeInUs(),
		res.GetOptimizeTimeInUs())

	if p.fd != nil {
		fmt.Fprintln(p.fd, s)
	}

	p.closeFile()
}
