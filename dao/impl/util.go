package dao_impl

import (
	"fmt"
	"learn-go-database/model"
	"reflect"
	"strconv"
	"strings"
	"time"
)

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

func GetColumns(user *model.User, columnFilter FilterColumn, columnList ...string) string {
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
func GetDataInsert(user *model.User) (string, string, []interface{}) {
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
func GetDataScan(user *model.User, include string) []interface{} {
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
