# About
helper program to get schema information on a SqlServer db.
To help power autocompletion feautres.

# pre set up
```bash
go install github.com/microsoft/go-mssqldb@latest
go get github.com/microsoft/go-mssqldb
```

# how to build

ideally use make
```bash
make -k
```

or if you are on windows with no make command just use the Go tooling directly
```bash
go build
```

or use -ldflags to omit symbol table, debug info, and dwarf symbol table. (smaller binary).
```bash
go build -ldflags="-s -w"
```

or disable bounds checks too
```bash
go build -gcflags=-B -ldflags="-s -w"
```