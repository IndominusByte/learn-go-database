package dao_impl

import (
	"context"
	"database/sql"
	"errors"
	"learn-go-database/dao"
	"learn-go-database/model"
	"strconv"
)

type userDaoImpl struct {
	DB *sql.DB
}

func CreateUserDaoImpl(db *sql.DB) dao.UserDao {
	return &userDaoImpl{DB: db}
}

func includeFilter(keyword string, data ...string) bool {
	for _, value := range data {
		if keyword == value {
			return true
		}
	}
	return false
}

func excludeFilter(keyword string, data ...string) bool {
	for _, value := range data {
		if keyword == value {
			return false
		}
	}
	return true
}

func (dao *userDaoImpl) Insert(ctx context.Context, user *model.User) (int64, error) {
	columns, values, data := GetDataInsert(user)

	query := "INSERT INTO users(" + columns + ")" + "VALUES(" + values + ") RETURNING id"
	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var lastId int64

	if err := stmt.QueryRowContext(ctx, data...).Scan(&lastId); err != nil {
		return 0, err
	}

	return lastId, nil
}

func (dao *userDaoImpl) FindById(ctx context.Context, id int64, include bool, columnList ...string) (model.User, error) {
	var filterColumn func(keyword string, data ...string) bool

	resultUser := &model.User{}

	if include {
		filterColumn = includeFilter
	} else {
		filterColumn = excludeFilter
	}

	columnScan := GetColumns(resultUser, filterColumn, columnList...)

	query := "SELECT " + columnScan + " FROM users WHERE id = $1 LIMIT 1"
	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return *resultUser, err
	}
	defer stmt.Close()

	if err := stmt.QueryRowContext(ctx, id).Scan(GetDataScan(resultUser, columnScan)...); err != nil {
		return *resultUser, errors.New("Id " + strconv.Itoa(int(id)) + " not found!")
	}

	return *resultUser, nil
}

func (dao *userDaoImpl) GetAll(ctx context.Context, include bool, columnList ...string) ([]model.User, error) {
	var filterColumn func(keyword string, data ...string) bool

	if include {
		filterColumn = includeFilter
	} else {
		filterColumn = excludeFilter
	}

	columnScan := GetColumns(&model.User{}, filterColumn, columnList...)
	query := "SELECT " + columnScan + " FROM users"
	stmt, err := dao.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		user       model.User
		resultUser []model.User
	)

	for rows.Next() {
		if err := rows.Scan(GetDataScan(&user, columnScan)...); err != nil {
			return nil, err
		}
		resultUser = append(resultUser, user)
	}

	return resultUser, nil
}
