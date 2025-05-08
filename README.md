# About
helper program to get schema information on a SqlServer db.
To help power autocompletion feautres.

# pre set up
```bash
go install github.com/microsoft/go-mssqldb@latest
go get github.com/microsoft/go-mssqldb
```

# build
```bash
go build
```

or use -ldflags to omit symbol table, debug info, and dwarf symbol table. (smaller binary).
```bash
go build -ldflags="-s -w"
```
