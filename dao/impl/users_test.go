package dao_impl

import (
	"context"
	"database/sql"
	"fmt"
	learngodatabase "learn-go-database"
	"learn-go-database/model"
	"testing"
)

func TestInsert(t *testing.T) {
	ctx := context.Background()
	db := CreateUserDaoImpl(learngodatabase.Connection())

	user1 := model.User{
		// Id:        sql.NullInt64{Int64: 200, Valid: true},
		Email:    sql.NullString{String: "oman5@gmail.com", Valid: true},
		Username: sql.NullString{String: "oman5", Valid: true},
		Password: sql.NullString{String: "asdasd", Valid: true},
		Balance:  sql.NullInt32{Int32: 20, Valid: true},
		Rating:   sql.NullFloat64{Float64: 20.2, Valid: true},
		// BirthDate: sql.NullTime{Time: time.Date(1999, 05, 19, 0, 0, 0, 0, time.Local), Valid: true},
		// Married:   sql.NullBool{Bool: true, Valid: true},
	}

	lastId, err := db.Insert(ctx, &user1)
	if err != nil {
		panic(err)
	}

	fmt.Println(lastId)
}

func TestFindById(t *testing.T) {
	ctx := context.Background()
	db := CreateUserDaoImpl(learngodatabase.Connection())

	var id int64 = 2

	value, err := db.FindById(ctx, id, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(value)
}

func TestGetAll(t *testing.T) {
	ctx := context.Background()
	db := CreateUserDaoImpl(learngodatabase.Connection())

	values, err := db.GetAll(ctx, false)
	if err != nil {
		panic(err)
	}

	fmt.Println(values)
}
