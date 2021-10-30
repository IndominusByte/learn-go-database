package learngodatabase

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
)

func TestExecSql(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	ctx := context.Background()
	_, err := conn.ExecContext(ctx, `INSERT INTO users(email, username, password, balance, rating, birth_date, married) 
		VALUES ('oman@gmail.com', 'oman', 'oman', 20, 2.4, '1999-9-19', true)`)
	if err != nil {
		panic(err)
	}

	fmt.Println("success insert new users")
}

func TestQuerySql(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	ctx := context.Background()
	rows, err := conn.QueryContext(ctx, "SELECT id, email, username, balance, rating, created_at, birth_date, married FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	type Users struct {
		Id                   int64
		Email                string
		Username             string
		Balance              int32
		Rating               sql.NullFloat64
		CreatedAt, BirthDate sql.NullTime
		Married              bool
	}

	var results []Users

	for rows.Next() {
		var email, username string
		var id int64
		var balance int32
		var rating sql.NullFloat64
		var created_at, birth_date sql.NullTime
		var married bool

		if err := rows.Scan(&id, &email, &username, &balance, &rating, &created_at, &birth_date, &married); err != nil {
			panic(err)
		}
		results = append(results, Users{
			Id:        id,
			Email:     email,
			Username:  username,
			Balance:   balance,
			Rating:    rating,
			CreatedAt: created_at,
			BirthDate: birth_date,
			Married:   married,
		})
	}

	fmt.Println(results)
}

func TestSqlInjection(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	username := "' or 1 = 1 -- -"
	password := "oman"

	ctx := context.Background()
	query := "SELECT * FROM users WHERE username = '" + username + "' AND password = '" + password + "' LIMIT 1"
	row, err := conn.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	if row.Next() {
		fmt.Println("LOGIN BERHASIL")
	} else {
		fmt.Println("LOGIN GAGAL")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	username := "' or 1 = 1 -- -"
	password := "oman"

	ctx := context.Background()
	query := "SELECT * FROM users WHERE username = $1 AND password = $2 LIMIT 1"
	row, err := conn.QueryContext(ctx, query, username, password)
	if err != nil {
		panic(err)
	}

	if row.Next() {
		fmt.Println("LOGIN BERHASIL")
	} else {
		fmt.Println("LOGIN GAGAL")
	}
}

func TestAutoIncrement(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	var lastId int
	msg := "Test"

	ctx := context.Background()
	query := "INSERT INTO comments(msg) VALUES($1) RETURNING id"
	err := conn.QueryRowContext(ctx, query, msg).Scan(&lastId)
	if err != nil {
		panic(err)
	}

	fmt.Println("Id", lastId)
}

func TestPreapreStatement(t *testing.T) {
	conn := Connection()
	defer conn.Close()

	ctx := context.Background()
	query := "INSERT INTO comments(msg) VALUES($1) RETURNING id"
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		var (
			lastId int
			msg    = "komen ke" + strconv.Itoa(i)
		)
		err := stmt.QueryRowContext(ctx, msg).Scan(&lastId)
		if err != nil {
			panic(err)
		}

		fmt.Println("id", lastId)
	}
}
