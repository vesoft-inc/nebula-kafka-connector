/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package printer

import (
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type DataSetPrinter struct {
	writer   table.Writer
	fd       *os.File
	filename string
}

func NewDataSetPrinter() DataSetPrinter {
	writer := table.NewWriter()
	configTableWriter(&writer, false)
	return DataSetPrinter{
		writer: writer,
	}
}

func (p *DataSetPrinter) ExportCsv(filename string) {
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

func (p *DataSetPrinter) closeFile() {
	if p.fd != nil {
		s := strings.Replace(p.writer.RenderCSV(), "\\\"", "", -1)
		fmt.Fprintln(p.fd, s)

		if err := p.fd.Close(); err != nil {
			fmt.Printf("Close file %s failed, %s", p.filename, err.Error())
		}
		p.fd = nil
		p.filename = ""
	}
}

func (p *DataSetPrinter) Print(res nebula.Result) {
	p.writer.ResetHeaders()
	p.writer.ResetRows()
	var header []interface{}
	for _, columName := range res.Columns() {
		header = append(header, columName)
	}
	p.writer.AppendHeader(table.Row(header))
	numCols := len(res.Columns())
	for res.HasNext() {
		var newRow []interface{}
		record, err := res.Next()
		if err != nil {
			continue
		}
		for j := 0; j < numCols; j++ {
			val, err := record.GetValueByIndex(j)
			if err != nil {
				continue
			}
			newRow = append(newRow, val.String())
		}
		p.writer.AppendRow(table.Row(newRow))
	}

	fmt.Println(p.writer.Render())
	p.closeFile()
}
