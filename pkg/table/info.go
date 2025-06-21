package table

import (
	"tradovatedataimport/pkg/csvdata"
	"tradovatedataimport/pkg/db"
)

type Info struct {
	name       string
	csvColumns csvdata.ColumnCollection
	dbColumns  db.ColumnCollection
}

type Column struct {
	InputColumn *csvdata.Column
	DbColumn    *db.Column
}

func NewInfo(name string, columns ...Column) *Info {
	numColumns := len(columns)
	result := &Info{
		name:       name,
		csvColumns: make([]*csvdata.Column, numColumns, numColumns),
		dbColumns:  make([]*db.Column, numColumns, numColumns),
	}
	for i, col := range columns {
		result.csvColumns[i] = col.InputColumn
		result.dbColumns[i] = col.DbColumn
	}

	return result
}

func (i *Info) Name() string {
	return i.name
}

func (i *Info) CsvColumns() csvdata.ColumnCollection {
	return i.csvColumns
}

func (i *Info) DbColumns() db.ColumnCollection {
	return i.dbColumns
}
