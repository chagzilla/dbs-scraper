package main

import (
	"os"

	"github.com/chagzilla/dbs-scraper/dbminer"
	"github.com/chagzilla/dbs-scraper/mongoMiner"
	"github.com/chagzilla/dbs-scraper/mySQLMiner"
)

type dbParser interface {
	PrintRecords(table string) error
}

func main() {

	var connection interface{}
	switch os.Args[2] {
	case "mongo":
		mongo, err := mongoMiner.New(os.Args[3], os.Args[4])
		if err != nil {
			panic(err)
		}
		connection = mongo
	case "mySQL":
		mysql, err := mySQLMiner.New(os.Args[3], os.Args[4])
		if err != nil {
			panic(err)
		}
		connection = mysql
	}
	switch os.Args[1] {
	case "records":
		printer := connection.(dbParser)
		if err := printer.PrintRecords(os.Args[5]); err != nil {
			panic(err)
		}
	case "schema":
		schemaMiner := connection.(dbminer.DatabaseMiner)
		if err := dbminer.Search(schemaMiner); err != nil {
			panic(err)
		}
	}
	// mm, err := mongoMiner.New(os.Args[2])
	// if err != nil {
	// 	panic(err)
	// }

	// if err := dbminer.Search(mm); err != nil {
	//	panic(err)
	// }
}
