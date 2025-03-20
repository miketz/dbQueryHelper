package main

import (
	"fmt"
)

func main() {
	fmt.Printf("test\n")
}

const sqlGetSchemas string = `select s.SCHEMA_NAME
from INFORMATION_SCHEMA.SCHEMATA s
order by s.SCHEMA_NAME`

func getSchemas() []string {
	return nil
}
