package repo_impl

import (
	"context"
	"database/sql"
	"fmt"
	"mygomod/db"
	"mygomod/error_my"
	"mygomod/model"
	"mygomod/model/req"
	"mygomod/mylog"
	"mygomod/repo"
	"time"

	"github.com/lib/pq"
)

type UserRepoImpl struct {
	sql *db.Sql
}
func NewUserRepo(sql *db.Sql) repo.UserRepo {
	return &UserRepoImpl{sql: sql}
}

// UpdateUser implements repo.UserRepo.
func (u *UserRepoImpl) UpdateUser(c context.Context, user model.User) (model.User, error) {
	statement := `
		UPDATE users 
		SET 
			full_name = (CASE WHEN LENGTH(:full_name) = 0 THEN full_name ELSE :full_name END)	,
			email = (CASE WHEN LENGTH(:email)=0 THEN email ELSE :email END)	,
			updated_at =COALESCE(:updated_at,updated_at)
		WHERE
			user_id=:user_id
	`
	user.UpdateAt = time.Now()
	result, err := u.sql.Db.NamedExecContext(c, statement, user)
	if err != nil {
		mylog.LogError(err)
		return user, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		mylog.LogError(err)
		return user, error_my.UserUpdateFail
	}
	if count == 0 {
		return user, error_my.UserUpdateFail
	}
	return user, nil

}

// SelectUserById implements repo.UserRepo.
func (u *UserRepoImpl) SelectUserById(c context.Context, userId string) (model.User, error) {
	user := model.User{}
	err := u.sql.Db.GetContext(c, &user, "SELECT * FROM users WHERE user_id=$1", userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, error_my.UserNotFound
		}
		fmt.Print(err)
		return user, err
	}
	return user, nil
}


func (u *UserRepoImpl) SaveUser(context context.Context, user model.User) (model.User, error) {
	statement := `
		INSERT INTO 
			users(
				user_id,
				email,
				password,
				role,
				full_name,
				created_at,
				updated_at
			)
			VALUES(
				:user_id,
				:email,
				:password,
				:role,
				:full_name,
				:created_at,
				:updated_at
			)
	`
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()
	_, err := u.sql.Db.NamedExecContext(context, statement, user)
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return user, error_my.UserConflict
			}
		}
	}
	return user, nil
}
func (u *UserRepoImpl) CheckLogIn(c context.Context, loginReq req.ReqLogIn) (model.User, error) {

	user := model.User{}
	err := u.sql.Db.GetContext(c, &user, "SELECT * FROM users WHERE email=$1", loginReq.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, error_my.UserNotFound

		}
		return user, err

	}
	return user, nil
}
