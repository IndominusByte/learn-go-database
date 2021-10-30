package dao

import (
	"context"
	"learn-go-database/model"
)

type UserDao interface {
	Insert(ctx context.Context, user *model.User) (int64, error)
	FindById(ctx context.Context, id int64, include bool, columnList ...string) (model.User, error)
	GetAll(ctx context.Context, include bool, columnList ...string) ([]model.User, error)
}
