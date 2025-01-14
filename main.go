package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/goombaio/namegenerator"
)

type Customer struct {
	ID   int
	Name string
	Tier int
}

var db *sql.DB
var nameGenerator namegenerator.Generator

func main() {

	seed := time.Now().UTC().UnixNano()
	nameGenerator = namegenerator.NewNameGenerator(seed)

	// Capture connection properties.
	cfg := mysql.Config{
		User:   "appuser",
		Passwd: "passme",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "forex",
	}
	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// result, err := db.Exec("INSERT INTO customer(name, tier) VALUES (@name, @tier)", sql.Named("name", "Peter"), sql.Named("tier", 1))
	result, err := db.Exec("INSERT INTO customer(name, tier) VALUES (?, ?)", nameGenerator.Generate(), 1)
	if err != nil {
		fmt.Errorf("insert customer: %v", err)
	}

	count, err := result.RowsAffected()
	if err != nil {
		fmt.Errorf("rows affected error: %v", err)
	}
	fmt.Println("Rows affected: ", count)

	customers := retrieveCustomers()
	fmt.Printf("customers: %+v\n", customers)
}

func retrieveCustomers() []Customer {
	rows, err := db.Query("SELECT id, name, tier FROM customer")
	if err != nil {
		log.Panic("db query error ", err)
	}
	defer rows.Close()

	var customers []Customer

	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.Tier); err != nil {
			log.Panic("db scan fail ", err)
		}
		customers = append(customers, customer)
	}
	if err := rows.Err(); err != nil {
		log.Panic("rows error ", err)
	}

	return customers
}
