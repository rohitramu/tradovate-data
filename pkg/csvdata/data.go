package csvdata

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"iter"
	"os"
	"strings"
	"tradovatedataimport/pkg/funcs"
)

type cleanFunc func(string) (string, error)

type ColumnCollection []*Column
type Row struct {
	columns ColumnCollection
	data    []string
}

type Column struct {
	name       string
	cleanFuncs []cleanFunc
}

func NewColumn(name string, cleanFuncs ...cleanFunc) *Column {
	return &Column{
		name:       name,
		cleanFuncs: cleanFuncs,
	}
}

func (c ColumnCollection) Rows(filepathCsv string) iter.Seq2[*Row, error] {
	return func(yield func(*Row, error) bool) {
		file, err := os.Open(filepathCsv)
		if err != nil {
			yield(nil, fmt.Errorf("failed to read %q: %v", filepathCsv, err))
			return
		}
		defer file.Close()

		reader := csv.NewReader(file)
		for rowNum := 0; ; rowNum++ {
			rowData, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				yield(nil, fmt.Errorf("failed to read row #%d: %v", rowNum, err))
				return
			}

			row := &Row{
				columns: c,
				data:    rowData,
			}

			// Validate header row.
			if rowNum == 0 {
				if err := row.validateHeader(); err != nil {
					yield(nil, fmt.Errorf("invalid header: %v", err))
					return
				}
				continue
			}

			// Return this row.
			if !yield(row, nil) {
				return
			}
		}
	}
}

func (r *Row) validateHeader() error {
	if len(r.data) != len(r.columns) {
		return fmt.Errorf("header has %d columns, but expected %d columns", len(r.data), len(r.columns))
	}

	for colNum, colInfo := range r.columns {
		if r.data[colNum] != colInfo.name {
			return fmt.Errorf("header at column %d was %q but expected %q", colNum, r.data[colNum], colInfo.name)
		}
	}

	return nil
}

func (r *Row) Clean() ([]string, error) {
	cleanedRow := make([]string, len(r.columns))
	for colNum, colInfo := range r.columns {
		val := r.data[colNum]

		// Run each of the cleanup functions in order.
		for _, cleanFunc := range colInfo.cleanFuncs {
			if cleanFunc == nil {
				cleanFunc = funcs.CleanNoOp
			}

			var err error
			val, err = cleanFunc(val)
			if err != nil {
				return nil, fmt.Errorf("failed to clean the value at column #%d: %v", colNum, err)
			}
		}

		cleanedRow[colNum] = strings.TrimSpace(val)
	}

	return cleanedRow, nil
}
