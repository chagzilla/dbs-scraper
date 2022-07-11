package mySQLMiner

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/chagzilla/dbs-scraper/dbminer"
	_ "github.com/go-sql-driver/mysql"
)

var sqlQuery = `SELECT TABLE_SCHEMA, TABLE_NAME, COLUMN_NAME FROM columns WHERE TABLE_SCHEMA NOT IN ('mysql', 'information_schema', 'performance_schema', 'sys') ORDER BY TABLE_SCHEMA, TABLE_NAME`

type MySQLMiner struct {
	Host     string
	Db       sql.DB
	database string
}

func New(host, db string) (*MySQLMiner, error) {
	m := MySQLMiner{Host: host, database: db}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MySQLMiner) connect() error {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf("root:password@tcp(%s:3306)/%s", m.Host, m.database))
	if err != nil {
		log.Panicln(err)
	}
	m.Db = *db
	return nil
}

func (m *MySQLMiner) PrintRecords(table string) error {
	rows, err := m.Db.Query("SELECT * FROM transactions")
	if err != nil {
		return err
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 1, ' ', 0)
	defer w.Flush()

	sep := []byte("\t")
	newLine := []byte("\n")

	w.Write([]byte(strings.Join(cols, "\t") + "\n"))

	row := make([][]byte, len(cols))
	rowPtr := make([]interface{}, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}

	for rows.Next() {
		_ = rows.Scan(rowPtr...)

		w.Write(bytes.Join(row, sep))
		w.Write(newLine)
	}
	return nil
}

func (m *MySQLMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema)

	schemarows, err := m.Db.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer schemarows.Close()

	var prevschema, prevtable string
	var db dbminer.Database
	var table dbminer.Table

	for schemarows.Next() {
		var currschema, currtable, currcol string
		if err := schemarows.Scan(&currschema, &currtable, &currcol); err != nil {
			return nil, err
		}

		if currschema != prevschema {

			if prevschema != "" {
				db.Tables = append(db.Tables, table)
				s.Databases = append(s.Databases, db)
			}
			db = dbminer.Database{Name: currschema, Tables: []dbminer.Table{}}
			prevschema = currschema
			prevtable = ""
		}

		if currtable != prevtable {

			if prevtable != "" {
				db.Tables = append(db.Tables, table)
			}
			table = dbminer.Table{Name: currtable, Columns: []string{}}
			prevtable = currtable
		}

		table.Columns = append(table.Columns, currcol)
	}
	db.Tables = append(db.Tables, table)
	s.Databases = append(s.Databases, db)
	if err := schemarows.Err(); err != nil {
		return nil, err
	}

	return s, nil
}
