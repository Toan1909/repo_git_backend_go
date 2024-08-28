package repo

import (
	"context"
	"mygomod/model"
	"mygomod/model/req"
)

type UserRepo interface {
	CheckLogIn(c context.Context, loginReq req.ReqLogIn) (model.User, error)
	SaveUser(c context.Context, user model.User) (model.User, error)
	SelectUserById(c context.Context, userId string) (model.User, error)
	UpdateUser(c context.Context, user model.User) (model.User, error)
}
