# goaddr

A Restful API (and db) for a super simple Address Book written in Golang.  Requires Go 1.16 or higher to be installed.

The following commands will download the executable into your go files and run the API server on port 8080:
``` bash
go get github.com/deversmann/goaddr
goaddr
```

The server will generate the db file the first time in the directory where you execute the command.  Subsequent executions in the same directory will reuse the same db file
