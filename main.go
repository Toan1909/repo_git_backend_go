package main

import (
	"fmt"
	"log"
	"mygomod/db"
	handler "mygomod/handler"
	"mygomod/helper"
	"mygomod/repo/repo_impl"
	"mygomod/router"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	sql := &db.Sql{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "190901",
		Dbname:   "golang01",
	}
	sql.ConnectDb()
	defer sql.CloseDb() //đóng db sau khi main kết thúc

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))
	userHandler := handler.UserHandler{
		UserRepo: repo_impl.NewUserRepo(sql),
	}
	repoGitHandler := handler.RepoGitHandler{
		GithubRepo: repo_impl.NewGitRepoImpl(sql),
	}
	api := router.API{
		Echo:        e,
		UserHandler: userHandler,
		RepoGitHandler: repoGitHandler,
	}
	api.SetupRouter()

	go scheduleUpdateTrending(30* time.Minute,repoGitHandler)

	e.Logger.Fatal(e.Start(":3000"))
}

func scheduleUpdateTrending (timeSchedule time.Duration, handler handler.RepoGitHandler) {
	ticker := time.NewTicker(timeSchedule)
	go func() {
		for {
			select {
			case <- ticker.C:
				fmt.Println( "Checking from github...")
				_,err:=helper.CrawTrendingRepos(handler.GithubRepo)
				if err!=nil {
					log.Println("ERR :",err.Error())
				}
			}
		}
	}()
}