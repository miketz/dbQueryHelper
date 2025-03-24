package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

// fake credentails
const connStr = "sqlserver://tester123:tester123@localhost/MSSQLSERVER01?database=OSHE_WIRS&TrustServerCertificate=true&Integrated Security=true&trusted_connection=yes"

func printCommands() {
	fmt.Printf(`Enter a command:
	schemas
	tables
	views
	cols
`)
}

func main() {
	if len(os.Args) < 2 { // GUARD: command line arg required
		printCommands()
		return
	}
	switch command := os.Args[1]; strings.ToLower(command) {
	case "schemas":
		printSchemas()
	case "tables":
		printTables()
	case "views":
		printViews()
	case "cols":
		printCols()
	default:
		printCommands()
	}
}

const SEP = "|"
const SEP_OUTER = ","

// print schemas in CSV format:
// schema|...
func printSchemas() {
	schemas := getSchemas()
	// print all but last with trailing SEP
	for i := 0; i < len(schemas)-1; i++ {
		fmt.Printf("%s%s", schemas[i], SEP)
	}
	// no SEP after last
	iLast := len(schemas) - 1
	fmt.Printf("%s", schemas[iLast])
}

// print tables in CSV format:
// schema|table, ...
func printTables() {
	tables := getTables()
	// print all but last with trailing SEP
	for i := 0; i < len(tables)-1; i++ {
		tab := tables[i]
		fmt.Printf("%s%s%s%s", tab.Schema, SEP, tab.Name, SEP_OUTER)
	}
	// no SEP after last
	iLast := len(tables) - 1
	last := tables[iLast]
	fmt.Printf("%s%s%s", last.Schema, SEP, last.Name)
}

// print views in CSV format:
// schema|view, ...
func printViews() {
	views := getViews()
	// print all but last with trailing SEP
	for i := 0; i < len(views)-1; i++ {
		view := views[i]
		fmt.Printf("%s%s%s%s", view.Schema, SEP, view.Name, SEP_OUTER)
	}
	// no SEP after last
	iLast := len(views) - 1
	last := views[iLast]
	fmt.Printf("%s%s%s", last.Schema, SEP, last.Name)
}

// print col data in CSV format:
// schema|table|col|dataType|maxLen|ordPos, ...
func printCols() {
	cols := getColData()
	// print all but last with trailing SEP
	for i := 0; i < len(cols)-1; i++ {
		col := cols[i]
		maxLen := ""
		if col.CharacterMaxLen.Valid {
			maxLen = strconv.Itoa(int(col.CharacterMaxLen.Int32))
		}
		fmt.Printf("%s%s%s%s%s%s%s%s%s%s%d%s",
			col.TableSchema, SEP,
			col.TableName, SEP,
			col.ColName, SEP,
			col.DataType, SEP,
			maxLen, SEP,
			col.OrdinalPositionn, SEP_OUTER)
	}
	// no SEP after last
	iLast := len(cols) - 1
	last := cols[iLast]
	maxLen := ""
	if last.CharacterMaxLen.Valid {
		maxLen = strconv.Itoa(int(last.CharacterMaxLen.Int32))
	}
	fmt.Printf("%s%s%s%s%s%s%s%s%s%s%d",
		last.TableSchema, SEP,
		last.TableName, SEP,
		last.ColName, SEP,
		last.DataType, SEP,
		maxLen, SEP,
		last.OrdinalPositionn)
}

// get relevant scheam names for autocompletion.
const sqlGetSchemas = `select s.SCHEMA_NAME
from INFORMATION_SCHEMA.SCHEMATA s
where s.SCHEMA_NAME not in (
	'db_accessadmin',
	'db_backupoperator',
	'db_datareader',
	'db_datawriter',
	'db_ddladmin',
	'db_denydatareader',
	'db_denydatawriter',
	'db_owner',
	'db_securityadmin'
)
order by s.SCHEMA_NAME`

func getSchemas() []string {
	schemas := make([]string, 0, 32)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(sqlGetSchemas)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var schema string
		err := rows.Scan(&schema)
		if err != nil {
			log.Fatal(err)
		}
		schemas = append(schemas, schema)
	}
	return schemas
}

type Table struct {
	Schema string
	Name   string
}

const sqlGetTables = `select t.TABLE_SCHEMA, t.TABLE_NAME
from INFORMATION_SCHEMA.TABLES t
where t.TABLE_TYPE='BASE TABLE'
order by t.TABLE_SCHEMA, t.TABLE_NAME`

func getTables() []Table {
	tables := make([]Table, 0, 1024)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(sqlGetTables)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		tab := Table{}
		err := rows.Scan(&tab.Schema, &tab.Name)
		if err != nil {
			log.Fatal(err)
		}
		tables = append(tables, tab)
	}
	return tables
}

type View struct {
	Schema string
	Name   string
}

const sqlGetViews = `select t.TABLE_SCHEMA, t.TABLE_NAME
from INFORMATION_SCHEMA.TABLES t
where t.TABLE_TYPE='VIEW'
order by t.TABLE_SCHEMA, t.TABLE_NAME`

func getViews() []View {
	views := make([]View, 0, 256)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(sqlGetViews)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		view := View{}
		err := rows.Scan(&view.Schema, &view.Name)
		if err != nil {
			log.Fatal(err)
		}
		views = append(views, view)
	}
	return views
}

type ColData struct {
	TableSchema      string
	TableName        string
	ColName          string
	DataType         string
	CharacterMaxLen  sql.NullInt32
	OrdinalPositionn int
}

const sqlGetColData = `select c.TABLE_SCHEMA, c.TABLE_NAME, c.COLUMN_NAME, c.DATA_TYPE, c.CHARACTER_MAXIMUM_LENGTH, c.ORDINAL_POSITION
from INFORMATION_SCHEMA.COLUMNS c
order by c.TABLE_SCHEMA, c.TABLE_NAME, c.ORDINAL_POSITION`

func getColData() []ColData {
	cols := make([]ColData, 0, 16384)

	db, err := sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query(sqlGetColData)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		col := ColData{}
		err := rows.Scan(&col.TableSchema, &col.TableName, &col.ColName, &col.DataType, &col.CharacterMaxLen, &col.OrdinalPositionn)
		if err != nil {
			log.Fatal(err)
		}
		cols = append(cols, col)
	}
	return cols
}
