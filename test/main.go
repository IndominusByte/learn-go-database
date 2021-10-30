package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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

type FilterColumn func(keyword string, data ...string) bool

func InColumnWithSeperator(keyword string, data string, seperator string) bool {
	var dataSlice = strings.Split(data, seperator)
	for _, value := range dataSlice {
		if keyword == value {
			return true
		}
	}
	return false
}

func ConvertDataToString(data interface{}) string {
	switch value := data.(type) {
	case string:
		return value
	case int32, int64:
		return fmt.Sprintf("%d", data)
	case float32, float64:
		return fmt.Sprintf("%.2f", data)
	case time.Time:
		return data.(time.Time).Format("2006-01-02 15:04:05")
	case bool:
		return strconv.FormatBool(data.(bool))
	default:
		panic("Unknown type data")
	}
}

func GetColumns(user *User, columnFilter FilterColumn, columnList ...string) string {
	typeof := reflect.TypeOf(*user)

	var resultColumn string

	for i := 0; i < typeof.NumField(); i++ {
		column := typeof.Field(i).Tag.Get("column")
		if columnFilter(column, columnList...) {
			resultColumn += column + ","
		}
	}

	return resultColumn[:len(resultColumn)-1]
}

// return (column, column) ($1, $2) (data, data)
func GetDataInsert(user *User) (string, string, []interface{}) {
	typeof, valueof := reflect.TypeOf(*user), reflect.ValueOf(user).Elem()

	var (
		resultColumn string
		resultStmt   string
		resultData   []interface{}
		stmtCount    = 1
	)

	for i := 0; i < valueof.NumField(); i++ {
		column := typeof.Field(i).Tag.Get("column")
		kind, value := typeof.Field(i).Type.Kind(), valueof.Field(i).Interface()
		if kind == reflect.Struct && reflect.ValueOf(value).Field(1).Interface().(bool) == true {
			resultData = append(resultData, ConvertDataToString(reflect.ValueOf(value).Field(0).Interface()))
			resultColumn += column + ","
			resultStmt += "$" + strconv.Itoa(stmtCount) + ","
			stmtCount++
		}
	}

	return resultColumn[:len(resultColumn)-1], resultStmt[:len(resultStmt)-1], resultData
}

// return (ptr user, ptr user)
func GetDataScan(user *User, include string) []interface{} {
	typeof, valueof := reflect.TypeOf(*user), reflect.ValueOf(user).Elem()

	var resultAddr []interface{}

	for i := 0; i < valueof.NumField(); i++ {
		column := typeof.Field(i).Tag.Get("column")
		if InColumnWithSeperator(column, include, ",") {
			resultAddr = append(resultAddr, valueof.Field(i).Addr().Interface())
		}
	}

	return resultAddr
}

func main() {
	user1 := User{
		Id:       sql.NullInt64{Int64: 200, Valid: true},
		Email:    sql.NullString{String: "oman1@gmail.com", Valid: true},
		Username: sql.NullString{String: "oman1", Valid: true},
		Password: sql.NullString{String: "asdasd", Valid: true},
		Balance:  sql.NullInt32{Int32: 20, Valid: true},
		Rating:   sql.NullFloat64{Float64: 20.2, Valid: true},
		// BirthDate: sql.NullTime{Time: time.Date(1999, 05, 19, 0, 0, 0, 0, time.Local), Valid: true},
		// Married:   sql.NullBool{Bool: true, Valid: true},
	}

	column, stmt, data := GetDataInsert(&user1)
	fmt.Println(column)
	fmt.Println(stmt)
	fmt.Println(data)
	fmt.Println("=========")

	columnScan := GetColumns(&User{}, func(keyword string, data ...string) bool {
		for _, value := range data {
			if keyword == value {
				return true
			}
		}
		return false
	}, "married", "birth_date")

	ptrUser := &User{}
	ptrUserScan := GetDataScan(ptrUser, columnScan)

	fmt.Println(columnScan)
	fmt.Println(ptrUserScan...)
}
