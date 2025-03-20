package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

// fake credentails
const connStr = "sqlserver://tester123:tester123@localhost/MSSQLSERVER01?database=OSHE_WIRS&TrustServerCertificate=true&Integrated Security=true&trusted_connection=yes"

func main() {
	fmt.Printf("starting...\n")

	cols := getColData()
	for _, col := range cols {
		fmt.Println("col: ", col)
	}

	fmt.Printf("done\n")
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
