package handle

import (
	"mygomod/error_my"
	"mygomod/helper"
	"mygomod/model"
	"mygomod/model/req"
	"mygomod/repo"
	"net/http"

	validator "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type RepoGitHandler struct {
	GithubRepo repo.GithubRepo
}

func (r *RepoGitHandler) HandleSaveReposFromGithub(c echo.Context) error {
	repos, _ := helper.CrawTrendingRepos(r.GithubRepo)
	for _, repo := range repos {
		_, err := r.GithubRepo.SaveRepo(c.Request().Context(), repo)
		if err != nil {
			return c.JSON(http.StatusConflict, model.Response{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
				Data:       nil,
			})
		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Craw TrendingRepos thành công",
		Data:       repos,
	})
}
func (r *RepoGitHandler) HandleGetReposFromDb(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	repos, err := r.GithubRepo.SelectRepos(c.Request().Context(),claims.UserId, 20)
	if err != nil {
		if err == error_my.RepoNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Get TrendingRepos from PostgreDb thành công",
		Data:       repos,
	})
}
func (r *RepoGitHandler) HandleGetRepoByName(c echo.Context) error {
	req := req.ReqGetRepoByName{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}

	repo, err := r.GithubRepo.SelectRepoByName(c.Request().Context(), req.Name)
	if err != nil {
		if err == error_my.RepoNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Get TrendingRepos from PostgreDb Byname thành công",
		Data:       repo,
	})
}

func (g *RepoGitHandler) HandleUpdateRepo(c echo.Context) error {
	repos, _ := helper.CrawTrendingRepos(g.GithubRepo)
	for _, repo := range repos {
		_, err := g.GithubRepo.UpdateRepo(c.Request().Context(), repo)
		if err != nil {
			if err == error_my.RepoNotFound {
				return c.JSON(http.StatusNotFound, model.Response{
					StatusCode: http.StatusNotFound,
					Message:    err.Error(),
					Data:       nil,
				})
			}
			return c.JSON(http.StatusUnprocessableEntity, model.Response{
				StatusCode: http.StatusUnprocessableEntity,
				Message:    err.Error(),
				Data:       nil,
			})

		}
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Update Repo from Git to Db thành công",
		Data:       repos,
	})

}
func (g *RepoGitHandler) HandleSelectBookmarks(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	repos, err := g.GithubRepo.SelectAllBookMarks(c.Request().Context(), claims.UserId)
	if err != nil {
		if err == error_my.BookMarkNotFound {
			return c.JSON(http.StatusNotFound, model.Response{
				StatusCode: http.StatusNotFound,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Select Bookmarks thành công",
		Data:       repos,
	})

}
func (g *RepoGitHandler) HandleBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	bId, _ := uuid.NewUUID()
	err := g.GithubRepo.BookMark(c.Request().Context(), bId.String(), req.NameRepo, claims.UserId)
	if err != nil {
		if err == error_my.BookMarkConflict {
			return c.JSON(http.StatusConflict, model.Response{
				StatusCode: http.StatusConflict,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Bookmark Repo thành công",
		Data:       nil,
	})

}
func (g *RepoGitHandler) HandleDeleteBookmark(c echo.Context) error {
	req := req.ReqBookmark{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			StatusCode: http.StatusBadRequest,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*model.JwtCustomclaims)
	err := g.GithubRepo.DeleteBookMark(c.Request().Context(), req.NameRepo, claims.UserId)
	if err != nil {
		if err == error_my.DeleteBookMarkFail {
			return c.JSON(http.StatusInternalServerError, model.Response{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
				Data:       nil,
			})
		}
		return c.JSON(http.StatusUnprocessableEntity, model.Response{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    err.Error(),
			Data:       nil,
		})
	}
	return c.JSON(http.StatusOK, model.Response{
		StatusCode: http.StatusOK,
		Message:    "Delete Bookmark Repo thành công",
		Data:       nil,
	})
}