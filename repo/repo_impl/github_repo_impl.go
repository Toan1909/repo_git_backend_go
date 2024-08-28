package repo_impl

import (
	"context"
	"database/sql"
	"mygomod/db"
	"mygomod/error_my"
	"mygomod/model"
	"mygomod/mylog"
	"mygomod/repo"
	"time"

	"github.com/lib/pq"
)

type GithubRepoImpl struct {
	sql *db.Sql
}

// BookMark implements repo.GithubRepo.
func (g *GithubRepoImpl) BookMark(c context.Context, bid string, nameRepo string, userId string) error {
	statement :=`
		INSERT INTO 
			bookmarks(
				bid,
				user_id,
				repo_name,
				created_at,
				updated_at
			)
			VALUES(
				$1,
				$2,
				$3,
				$4,
				$5
			)
		`
	timeNow := time.Now()
	_,err := g.sql.Db.ExecContext(c,statement,bid,userId,nameRepo,timeNow,timeNow)
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return error_my.BookMarkConflict
			}
		}
		return err
	}
	return nil
}

// DeleteBookMark implements repo.GithubRepo.
func (g *GithubRepoImpl) DeleteBookMark(c context.Context, nameRepo string, uId string) error {
	result:= g.sql.Db.MustExecContext(c,`DELETE FROM bookmarks where user_id=$1 AND repo_name=$2`,uId,nameRepo)
	_,err:=result.RowsAffected()
	if err!=nil {
		mylog.LogError(err)
		return error_my.DeleteBookMarkFail
	}
	return  nil
}

// SelectAllBookMarks implements repo.GithubRepo.
func (g *GithubRepoImpl) SelectAllBookMarks(c context.Context, userId string) ([]model.GithubRepo, error) {
	statement:= `
		SELECT 
			repos.name,
			repos.description,
			repos.url,
			repos.color,
			repos.lang,
			repos.fork,
			repos.stars,
			repos.stars_today,
			repos.build_by,
			COALESCE(repos.name=bookmarks.repo_name,false) as bookmarked
		FROM repos 
		INNER JOIN bookmarks 
		ON repos.name=bookmarks.repo_name AND bookmarks.user_id=$1
		`
	var repos []model.GithubRepo;
	err := g.sql.Db.SelectContext(c, &repos, statement,userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return repos, error_my.RepoNotFound
		}
		return repos, err
	}
	return repos, nil
}

func NewGitRepoImpl(sql *db.Sql) repo.GithubRepo {
	return &GithubRepoImpl{sql: sql}
}

// SaveRepo implements repo.GithubRepo.
func (g *GithubRepoImpl) SaveRepo(c context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	statement := `
		INSERT INTO 
			repos(
				name,
				description,
				url,
				color,
				lang,
				fork,
				stars,
				stars_today,
				build_by,
				created_at,
				updated_at
			)
			VALUES(
				:name,
				:description,
				:url,
				:color,
				:lang,
				:fork,
				:stars,
				:stars_today,
				:build_by,
				:created_at,
				:updated_at	
			)
	`
	repo.Created_at = time.Now()
	repo.Updated_at = time.Now()
	_, err := g.sql.Db.NamedExecContext(c, statement, repo)
	if err != nil {
		mylog.LogError(err)
		if err, ok := err.(*pq.Error); ok {
			if err.Code.Name() == "unique_violation" {
				return repo, error_my.RepoConflict
			}
		}
	}
	return repo, nil
}

// SelectRepoByName implements repo.GithubRepo.
func (g *GithubRepoImpl) SelectRepoByName(c context.Context, name string) (model.GithubRepo, error) {
	
	var repo model.GithubRepo
	err := g.sql.Db.GetContext(c, &repo, `SELECT 
					repos.name,
					repos.description,
					repos.url,
					repos.color,
					repos.lang,
					repos.fork,
					repos.stars,
					repos.stars_today,
					repos.build_by
				FROM repos 
				WHERE name=$1`,name)
	if err != nil {
		if err == sql.ErrNoRows {
			return repo, error_my.RepoNotFound
		}
		return repo, err
	}
	return repo, nil
}

// SelectRepos implements repo.GithubRepo.
func (g *GithubRepoImpl) SelectRepos(c context.Context,userId string ,limit int) ([]model.GithubRepo, error) {
	var repos []model.GithubRepo
	statement:=`SELECT 
					repos.name,
					repos.description,
					repos.url,
					repos.color,
					repos.lang,
					repos.fork,
					repos.stars,
					repos.stars_today,
					repos.build_by,
					COALESCE(repos.name=bookmarks.repo_name,false) as bookmarked
				FROM repos 
				FULL OUTER JOIN bookmarks
				ON repos.name = bookmarks.repo_name AND bookmarks.user_id=$1 
				ORDER BY repos.updated_at desc LIMIT $2`
	err := g.sql.Db.SelectContext(c, &repos,statement,userId, limit)
	if err != nil {
		if err == sql.ErrNoRows {
			return repos, error_my.RepoNotFound
		}
		return repos, err
	}
	return repos, nil
}

// UpdateRepo implements repo.GithubRepo.
func (g *GithubRepoImpl) UpdateRepo(c context.Context, repo model.GithubRepo) (model.GithubRepo, error) {
	repo.Updated_at = time.Now()
	statement := `
		UPDATE repos 
		SET 
			description = (CASE WHEN LENGTH(:description)=0 THEN description ELSE :description END)	,
			url = (CASE WHEN LENGTH(:url)=0 THEN url ELSE :url END)	,
			color = (CASE WHEN LENGTH(:color)=0 THEN color ELSE :color END)	,
			lang = (CASE WHEN LENGTH(:lang)=0 THEN lang ELSE :lang END)	,
			fork = (CASE WHEN LENGTH(:fork)=0 THEN fork ELSE :fork END)	,
			stars = (CASE WHEN LENGTH(:stars)=0 THEN stars ELSE :stars END)	,
			stars_today = (CASE WHEN LENGTH(:stars_today)=0 THEN stars_today ELSE :stars_today END)	,
			build_by = (CASE WHEN LENGTH(:build_by)=0 THEN build_by ELSE :build_by END)	,
			updated_at =COALESCE(:updated_at,updated_at),
			created_at	=COALESCE(:created_at,created_at)
		WHERE
			name=:name
	`
	result, err := g.sql.Db.NamedExecContext(c, statement, repo)
	if err != nil {
		mylog.LogError(err)
		return repo, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		mylog.LogError(err)
		return repo, error_my.RepoUpdateFail
	}
	if count == 0 {
		return repo, error_my.RepoUpdateFail
	}
	return repo, nil
}
