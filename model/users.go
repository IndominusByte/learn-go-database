package model

import "database/sql"

type User struct {
	Id        sql.NullInt64   `column:"id"`
	Email     sql.NullString  `column:"email"`
	Username  sql.NullString  `column:"username"`
	Password  sql.NullString  `column:"password"`
	Balance   sql.NullInt32   `column:"balance"`
	Rating    sql.NullFloat64 `column:"rating"`
	CreatedAt sql.NullTime    `column:"created_at"`
	BirthDate sql.NullTime    `column:"birth_date"`
	Married   sql.NullBool    `column:"married"`
}
