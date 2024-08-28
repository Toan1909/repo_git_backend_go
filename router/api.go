package router

import (
	handle "mygomod/handler"
	"mygomod/middleware"

	"github.com/labstack/echo/v4"
)

type API struct {
	Echo           *echo.Echo
	UserHandler    handle.UserHandler
	RepoGitHandler handle.RepoGitHandler
}

func (api *API) SetupRouter() {
	// user
	api.Echo.POST("/user/sign-in", api.UserHandler.HandleSignIn)
	api.Echo.POST("/user/sign-up", api.UserHandler.HandleSignUp)
	// profile
	user := api.Echo.Group("/user", middleware.JWTMiddleWare())
	user.GET("/profile", api.UserHandler.HandleProfile)
	user.PUT("/profile/update", api.UserHandler.HandleUpdateUser)
	//repo github
	github := api.Echo.Group("/github", middleware.JWTMiddleWare())
	github.GET("/trending",api.RepoGitHandler.HandleGetReposFromDb)
	github.POST("/crawl", api.RepoGitHandler.HandleSaveReposFromGithub)
	github.POST("/get-repo",api.RepoGitHandler.HandleGetRepoByName)
	//bookmark
	bookmark :=api.Echo.Group("/bookmark",middleware.JWTMiddleWare())
	bookmark.GET("/list",api.RepoGitHandler.HandleSelectBookmarks)
	bookmark.POST("/add",api.RepoGitHandler.HandleBookmark)
	bookmark.DELETE("/delete",api.RepoGitHandler.HandleDeleteBookmark)

}
