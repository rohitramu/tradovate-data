package db

import (
	"fmt"
	"strings"
)

func (c ColumnCollection) CreateTableSql(tableName string) string {
	primaryKeys := []string{}
	tableColumns := strings.Builder{}
	for i, col := range c {
		if col.isPrimaryKey {
			primaryKeys = append(primaryKeys, col.name)
		}
		if i != 0 {
			tableColumns.WriteString(",\n")
		}
		tableColumns.WriteString(fmt.Sprintf("  %s %s", col.name, string(col.dataType)))
	}
	if len(primaryKeys) > 0 {
		tableColumns.WriteString(fmt.Sprintf(",\n  PRIMARY KEY(%s)", strings.Join(primaryKeys, ", ")))
	}

	return fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n%s) WITHOUT ROWID;", tableName, tableColumns.String())
}

func (c ColumnCollection) InsertRowsSql(tableName string, rows ...[]string) (string, error) {
	numRows := len(rows)
	if numRows == 0 {
		return "", nil
	}

	valuesSql := strings.Builder{}
	for i, row := range rows {
		rowStr := fmt.Sprintf("\"%s\"", strings.Join(row, "\", \""))
		if len(c) != len(row) {
			return "", fmt.Errorf("table %q has %d columns, but only %d values were provided: %s", tableName, len(c), len(row), rowStr)
		}

		valuesSql.WriteString(fmt.Sprintf("\n  (%s)", rowStr))

		// Don't add a newline for the final row.
		if i+1 != numRows {
			valuesSql.WriteRune(',')
		}
	}

	return fmt.Sprintf("INSERT OR REPLACE INTO %s VALUES %s;", tableName, valuesSql.String()), nil
}
