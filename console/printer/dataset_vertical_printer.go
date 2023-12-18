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
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Printf("Open or Create file %s failed, %s", filename, err.Error())
		return
	}
	p.fd = fd
	p.filename = filename
}

func (p *DataSetVerticalPrinter) PrintDataSet(res *nebula.ResultSet) {
	if res.GetColSize() == 0 {
		return
	}

	maxColumnWidth := 0
	for _, columName := range res.GetColNames() {
		if len(columName) > maxColumnWidth {
			maxColumnWidth = len(columName)
		}
	}
	numRows := res.GetRowSize()
	numCols := res.GetColSize()
	print := func(str string) {
		fmt.Println(str)
		if p.fd != nil {
			fmt.Fprintln(p.fd, str)
		}
	}
	for i := 0; i < numRows; i++ {
		print(fmt.Sprintf("*************************** %d. row ***************************", i+1))
		record, err := res.GetRowValuesByIndex(i)
		if err != nil {
			continue
		}
		for j := 0; j < numCols; j++ {
			val, err := record.GetValueByIndex(j)
			if err != nil {
				continue
			}
			colName := res.GetColNames()[j]
			print(fmt.Sprintf("%*s: %s", maxColumnWidth, colName, val.String()))
		}
	}

	if p.fd != nil {
		if err := p.fd.Close(); err != nil {
			fmt.Printf("Close file %s failed, %s", p.filename, err.Error())
		}
		p.fd = nil
		p.filename = ""
	}
}
