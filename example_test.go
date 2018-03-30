package goql_test

import (
	"database/sql"
	"fmt"
	"log"
)

func Example() {
	con, err := sql.Open("goql", "net/http")
	if err != nil {
		log.Fatal(err)
	}
	defer con.Close()

	rows, err := con.Query("SELECT name, receiver FROM funcs WHERE name='Do'") // client.Do
	if err != nil {
		log.Fatal(err)
	}
	var name, rec string
	for rows.Next() {
		if err := rows.Scan(&name, &rec); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s.%s", rec, name)
	}
	// Output:
	// Client.Do
}
