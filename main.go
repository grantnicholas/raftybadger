package main

import (
	"fmt"
	"raftybadger/badgerdb"
	"raftybadger/server"
)

func main() {
	fmt.Printf("...Starting up...\r\n")
	db := badgerdb.GetDB()
	server.Serve(db)
}
