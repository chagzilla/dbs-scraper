package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float32 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

func main() {
	session, err := mgo.Dial("10.0.0.197")
	if err != nil {
		log.Panicln(err)
	}
	defer session.Close()

	results := make([]Transaction, 0)
	if err := session.DB("store").C("transaction").Find(nil).All(&results); err != nil {
		log.Panicln(err)
	}
	for _, txn := range results {
		fmt.Println(txn.CCNum, txn.Date, txn.Amount, txn.Cvv, txn.Expiration)
	}
}
