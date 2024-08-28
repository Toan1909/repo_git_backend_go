package helper

import (
	"context"
	"fmt"
	"log"
	"mygomod/error_my"
	"mygomod/model"
	"mygomod/repo"
	"regexp"
	"runtime"
	"strings"

	"github.com/gocolly/colly/v2"
)

func CrawTrendingRepos(githubRepo repo.GithubRepo) ([]model.GithubRepo, error) {
	var repos []model.GithubRepo

	// Initialize the collector
	c := colly.NewCollector(
		colly.AllowedDomains("github.com"),
	)
	
	// Find and extract the repository details
	c.OnHTML("article.Box-row", func(e *colly.HTMLElement) {
		githubRepo := model.GithubRepo{}
		// filter 4 data repo
		n := strings.Replace(e.ChildText("h2.h3.lh-condensed a"), "\n", "", -1)
		githubRepo.Name = strings.Replace(n, " ", "", -1)
		githubRepo.Url = strings.TrimSpace(e.ChildAttr("h2.h3.lh-condensed a", "href"))
		githubRepo.Description = strings.TrimSpace(e.ChildText("p.col-9.color-fg-muted.my-1.pr-4"))
		bgColor := strings.TrimSpace(e.ChildAttr("span.d-inline-block.ml-0.mr-3 > span.repo-language-color", "style"))
		
		match := regexp.MustCompile("#[a-zA-z0-9_]+").FindStringSubmatch(bgColor)
		if len(match) > 0 {
			githubRepo.Color = match[0]
		}

		githubRepo.Lang = strings.TrimSpace(e.ChildText("span.d-inline-block.ml-0.mr-3 span[itemprop=programmingLanguage]"))
		e.ForEach("a.Link--muted.d-inline-block.mr-3", func(_ int, el *colly.HTMLElement) {
			if strings.Contains(el.Attr("href"), "/stargazer") {
				githubRepo.Stars = strings.TrimSpace(el.Text)
			}
		})
		e.ForEach("a.Link--muted.d-inline-block.mr-3", func(_ int, el *colly.HTMLElement) {
			if strings.Contains(el.Attr("href"), "/forks") {
				githubRepo.Fork = strings.TrimSpace(el.Text)
			}
		})
		var buildBy []string
		e.ForEach("a.d-inline-block img.avatar", func(_ int, el *colly.HTMLElement) {
			avatarContributor := el.Attr("src")
			buildBy = append(buildBy, avatarContributor)
		})
		githubRepo.Build_by = strings.Join(buildBy, ",")
		githubRepo.Stars_today = strings.TrimSpace(e.ChildText("span.d-inline-block.float-sm-right"))
		repos = append(repos, githubRepo)
		log.Println("Crawl thành công repo :",githubRepo.Name)
	})

	
	c.OnScraped(func(r *colly.Response) {
		log.Println("Start onScraped")
		queue := NewJobQueue(runtime.NumCPU())
		queue.Start()
		log.Println("Queue started with workers:", runtime.NumCPU())
		defer queue.Stop()
		for _, repo := range repos {
			fmt.Println("Submitting repo:", repo.Name)
			queue.Submit(&RepoProcess{
				Repo:       repo,
				GithubRepo: githubRepo,
			})
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	// Visit the GitHub trending page
	err := c.Visit("https://github.com/trending")
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// imple interface Job
type RepoProcess struct {
	Repo       model.GithubRepo
	GithubRepo repo.GithubRepo
}

func (rp *RepoProcess) Process() {
	log.Println("Processing repo:", rp.Repo.Name)
	//select repo by name
	cacheRepo, err := rp.GithubRepo.SelectRepoByName(context.Background(), rp.Repo.Name)
	if err != nil {
		if err == error_my.RepoNotFound {
			//Không tìm thấy repo này trong db -> tiến hành add repo vào db
			fmt.Println("ADD Repo :", rp.Repo.Name)
			rp.GithubRepo.SaveRepo(context.Background(), rp.Repo)
			return
		}
		log.Println(err.Error())
	}
	//Trường hợp tìm thấy repo đã có sẵn -> thì update
	if rp.Repo.Stars != cacheRepo.Stars ||
		rp.Repo.Fork != cacheRepo.Fork ||
		rp.Repo.Stars_today != cacheRepo.Stars_today {
		rp.Repo.Created_at = cacheRepo.Created_at
		fmt.Println("UPDATE Repo :", rp.Repo.Name)
		_, err := rp.GithubRepo.UpdateRepo(context.Background(), rp.Repo)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
