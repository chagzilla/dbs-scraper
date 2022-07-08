package main

import (
	"os"

	"github.com/chagzilla/dbs-scraper/dbminer"
	"github.com/chagzilla/dbs-scraper/mongoMiner"
)

func main() {
	mm, err := mongoMiner.New(os.Args[1])
	if err != nil {
		panic(err)
	}

	if err := dbminer.Search(mm); err != nil {
		panic(err)
	}
}
