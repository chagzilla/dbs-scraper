package dbminer

import (
	"fmt"
	"regexp"
)

type DatabaseMiner interface {
	GetSchema() (*Schema, error)
}

type Schema struct {
	Databases []Database
}

type Table struct {
	Name    string
	Columns []string
}

func Search(m DatabaseMiner) error {
	s, err := m.GetSchema()
	if err != nil {
		return err
	}

	re := getRegex()
	for _, database := range s.Databases {
		for _, table := range database.Tables {
			for _, field := range table.Columns {
				for _, r := range re {
					if r.MatchString(field) {
						fmt.Println(database)
						fmt.Printf("[+] HIT: %s\n", field)
					}
				}
			}
		}
	}
	return nil
}

func getRegex() []*regexp.Regexp {
	return nil
}

func main() {
	fmt.Println("vim-go")
}
