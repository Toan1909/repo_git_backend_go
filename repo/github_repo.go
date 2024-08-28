package repo
import(
	"context"
	"mygomod/model"
)
type GithubRepo interface {
	//repository
	SaveRepo(c context.Context, repo model.GithubRepo) (model.GithubRepo, error)
	SelectRepos(c context.Context, userId string,limit int) ([]model.GithubRepo, error)
	SelectRepoByName(c context.Context, name string) (model.GithubRepo, error)
	UpdateRepo(c context.Context, repo model.GithubRepo) (model.GithubRepo, error)
	//bookmark
	SelectAllBookMarks(c context.Context,userId string) ([]model.GithubRepo,error)
	BookMark(c context.Context,bid,nameRepo,userId string) error
	DeleteBookMark(c context.Context,nameRepo,uId string) error
}