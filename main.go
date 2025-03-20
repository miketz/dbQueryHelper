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

	views := getViews()
	for _, view := range views {
		fmt.Printf("%v\n", view)
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
