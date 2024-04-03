/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package printer

import (
	"fmt"
	"os"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type DataSetVerticalPrinter struct {
	fd       *os.File
	filename string
}

func NewDataSetVerticalPrinter() DataSetVerticalPrinter {
	return DataSetVerticalPrinter{}
}

func (p *DataSetVerticalPrinter) ExportCsv(filename string) {
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

func (p *DataSetVerticalPrinter) closeFile() {
	if p.fd != nil {
		if err := p.fd.Close(); err != nil {
			fmt.Printf("Close file %s failed, %s", p.filename, err.Error())
		}
		p.fd = nil
		p.filename = ""
	}
}

func (p *DataSetVerticalPrinter) Print(res nebula.Result) {
	if res.RowSize() == 0 {
		return
	}

	maxColumnWidth := 0
	for _, columName := range res.Columns() {
		if len(columName) > maxColumnWidth {
			maxColumnWidth = len(columName)
		}
	}
	numCols := len(res.Columns())
	print := func(str string) {
		fmt.Println(str)
		if p.fd != nil {
			fmt.Fprintln(p.fd, str)
		}
	}
	var i int
	for res.HasNext() {
		print(fmt.Sprintf("*************************** %d. row ***************************", i+1))
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
			print(fmt.Sprintf("%*s: %s", maxColumnWidth, colName, val.String()))
		}
	}

	p.closeFile()
}
