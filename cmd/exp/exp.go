package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/tepavcevic/go-template-server/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()

	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	fmt.Println("Connected!")

	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name TEXT,
	// 	email TEXT UNIQUE NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// 	`)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("tables created")

	// name := "Djoe"
	// email := "djoordje.tepavcevic@tob.ba"
	// row := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING 1;", name, email)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("inserted id:", id)

	// id := 1
	// row := db.QueryRow(`
	// 		SELECT name, email
	// 		FROM users
	// 		WHERE id=$1`,
	// 	id)
	// var name, email string
	// err = row.Scan(&name, &email)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User %s with %s\n", name, email)

	// userId := 1
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err = db.Exec(`
	// 	INSERT INTO orders (user_id, amount, description)
	// 	VALUES ($1, $2, $3);`, userId, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Added fake orders")
	// }

	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Description string
	}

	var orders []Order
	userID := 1
	rows, err := db.Query(`
	SELECT id, amount, description
	FROM orders
	WHERE user_id=$1
	`, userID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		order.UserID = userID
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		orders = append(orders, order)
	}
	if rows.Err() != nil {
		panic(err)
	}

	fmt.Println("Orders:", orders)
}
