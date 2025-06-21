package main

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
	"tradovatedataimport/pkg/table"

	_ "github.com/mattn/go-sqlite3"
)

const DEBUG = false

var debugLogger *log.Logger = log.New(io.Discard, "", 0)

const SQL_INSERT_BATCH_SIZE = 100000

var dataKindToTableInfo = map[string]*table.Info{
	"performance": table.Performance(),
	"cash":        table.Cash(),
}

func init() {
	if DEBUG {
		logFilepath := fmt.Sprintf("./import-%s.log", time.Now().Format("20060102150405"))
		logFile, err := os.Create(logFilepath)
		if err != nil {
			panic(fmt.Sprintf("Failed to create log file %q: %v", logFilepath, err))
		}
		debugLogger = log.New(logFile, "", log.Flags())
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %v <data_dir> <db_filepath>", os.Args[0])
	}

	dataDir := os.Args[1]
	dbFilepath := os.Args[2]

	if filepath.Ext(dbFilepath) != ".db" {
		log.Fatalf("Output filepath must have an extension of \".db\".")
	}

	if err := importDataDir(dataDir, dbFilepath); err != nil {
		log.Fatalf("failed to import data from %q: %v", dataDir, err)
	}
}

func importDataDir(dataDir string, dbFilepath string) error {
	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to visit path %q: %v", path, err)
		}

		// Skip any path that doesn't point to a CSV file.
		basename, isCsvFile := strings.CutSuffix(d.Name(), ".csv")
		if !isCsvFile {
			return nil
		}
		// Skip unknown tables.
		tableInfo, found := dataKindToTableInfo[basename]
		if !found {
			return nil
		}

		// Import data into DB.
		if err := importDataFile(tableInfo, path, dbFilepath); err != nil {
			return fmt.Errorf("failed to import data file %q: %v", path, err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed while traversing data directory: %v", err)
	}

	return nil
}

func importDataFile(tableInfo *table.Info, dataFilepath string, dbFilepath string) error {
	// Connect to the database and create it if it doesn't exist.
	dbConn, err := sql.Open("sqlite3", dbFilepath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer dbConn.Close()

	createTableStatement := tableInfo.DbColumns().CreateTableSql(tableInfo.Name()) + "\n"
	debugLogger.Println(createTableStatement)
	if _, err := dbConn.Exec(createTableStatement); err != nil {
		return fmt.Errorf("failed to execute \"CREATE TABLE\" statement: %v", err)
	}

	// Count the number of rows.
	numRows := 0
	for _, err := range tableInfo.CsvColumns().Rows(dataFilepath) {
		if err != nil {
			log.Fatalf("Failed to read file %q: %v", dataFilepath, err)
		}
		numRows++
	}

	// Clean and insert the data in the DB.
	rowNum := 0
	batchRows := make([][]string, 0, SQL_INSERT_BATCH_SIZE)
	for row, err := range tableInfo.CsvColumns().Rows(dataFilepath) {
		if err != nil {
			log.Fatalf("Failed to read file %q: %v", dataFilepath, err)
		}

		// Clean the data
		cleanedRow, err := row.Clean()
		if err != nil {
			return fmt.Errorf("failed to clean row #%d: %v", rowNum, err)
		}

		// Add row to the batch.
		batchRows = append(batchRows, cleanedRow)
		rowNum++

		// If we haven't hit the batch size and haven't reached the last row, don't insert the batch
		// into the DB yet.
		if rowNum%SQL_INSERT_BATCH_SIZE != 0 && rowNum != numRows {
			continue
		}

		// Generate the SQL statement to insert the row into the DB.
		rowsInsertStatement, err := tableInfo.DbColumns().InsertRowsSql(tableInfo.Name(), batchRows...)
		if err != nil {
			return fmt.Errorf("failed to generate \"INSERT\" statement: %v", err)
		}
		rowsInsertStatement += "\n"
		debugLogger.Println(rowsInsertStatement)
		// Insert the batch of rows into the DB.
		if DEBUG {
			startRowNum := rowNum - SQL_INSERT_BATCH_SIZE
			if startRowNum < 0 {
				startRowNum = 0
			}
			log.Printf("DEBUG: [%s] Importing rows %d-%d/%d from %s", tableInfo.Name(), startRowNum+1, rowNum, numRows, dataFilepath)
		}
		if _, err := dbConn.Exec(rowsInsertStatement); err != nil {
			return fmt.Errorf("failed to execute \"INSERT\" statement: %v", err)
		}

		// Start a new batch.
		batchRows = make([][]string, 0, SQL_INSERT_BATCH_SIZE)
	}

	log.Printf("Wrote %d rows to table %q: %s", numRows, tableInfo.Name(), dataFilepath)

	return nil
}
