package mongoMiner

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/chagzilla/dbs-scraper/dbminer"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float32 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

type MongoMiner struct {
	Host     string
	session  *mgo.Session
	database string
}

func New(host, db string) (*MongoMiner, error) {
	m := MongoMiner{Host: host, database: db}
	err := m.connect()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *MongoMiner) connect() error {
	s, err := mgo.Dial(m.Host)
	if err != nil {
		return err
	}
	m.session = s
	return nil
}

func (m *MongoMiner) PrintRecords(table string) error {

	results := make([]Transaction, 0)
	if err := m.session.DB(m.database).C(table).Find(nil).All(&results); err != nil {
		return err
	}
	for _, txn := range results {
		fmt.Println(txn.CCNum, txn.Date, txn.Amount, txn.Cvv, txn.Expiration)
	}
	return nil
}

func (m *MongoMiner) GetSchema() (*dbminer.Schema, error) {
	var s = new(dbminer.Schema)

	dbnames, err := m.session.DatabaseNames()
	if err != nil {
		return nil, err
	}

	for _, dbname := range dbnames {
		db := dbminer.Database{Name: dbname, Tables: []dbminer.Table{}}
		collections, err := m.session.DB(dbname).CollectionNames()
		if err != nil {
			return nil, err
		}

		for _, collection := range collections {
			table := dbminer.Table{Name: collection, Columns: []string{}}

			var docRaw bson.Raw
			err := m.session.DB(dbname).C(collection).Find(nil).One(&docRaw)
			if err != nil {
				return nil, err
			}

			var doc bson.RawD
			if err := docRaw.Unmarshal(&doc); err != nil {
				if err != nil {
					return nil, err
				}
			}

			for _, f := range doc {
				table.Columns = append(table.Columns, f.Name)
			}
			db.Tables = append(db.Tables, table)
		}
		s.Databases = append(s.Databases, db)
	}
	return s, nil
}
