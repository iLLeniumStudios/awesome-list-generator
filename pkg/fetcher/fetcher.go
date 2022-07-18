package fetcher

import (
	"context"
	"github.com/google/go-github/v45/github"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/config"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/models"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"os"
	"strings"
)

type Fetcher interface {
	Fetch() (models.AwesomeList, error)
}

type fetcher struct {
	Config *config.Config
	Client *github.Client
}

func MakeGithubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	return client
}

func New(conf *config.Config) Fetcher {
	return &fetcher{
		Config: conf,
		Client: MakeGithubClient(),
	}
}

func (f *fetcher) GetGithubReposForUsername(user string) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	opts := &github.RepositoryListOptions{ListOptions: github.ListOptions{PerPage: 100}}
	for {
		repos, resp, err := f.Client.Repositories.List(context.Background(), user, opts)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return allRepos, nil
}

func (f *fetcher) GetFiveMReposForUser(user config.User, ignoreList utils.StringList) ([]models.Repository, error) {
	var userRepositories []models.Repository
	repos, err := f.GetGithubReposForUsername(user.Name)
	if err != nil {
		return nil, err
	}

	for _, repo := range repos {
		if *repo.Private || *repo.StargazersCount < f.Config.MinStars || *repo.Disabled {
			continue
		}
		if ignoreList.Contains(*repo.Name) || *repo.Name == user.Name {
			continue
		}
		repoObject := models.Repository{
			Name: *repo.Name,
			URL:  *repo.HTMLURL,
		}
		if *repo.Fork && !user.SkipForks {
			if *repo.StargazersCount < f.Config.MinStarsForFork {
				continue
			}
			fullRepo, _, err := f.Client.Repositories.Get(context.Background(), user.Name, *repo.Name)
			if err != nil {
				continue
			}
			splitted := strings.Split(*fullRepo.Parent.Owner.HTMLURL, "/")
			comparison, _, err := f.Client.Repositories.CompareCommits(context.Background(), user.Name, *repo.Name, *fullRepo.DefaultBranch, splitted[len(splitted)-1]+":"+*fullRepo.Parent.DefaultBranch, nil)
			if err != nil {
				panic(err)
			}
			repoObject.Fork = &models.Fork{
				Owner:    *fullRepo.Parent.Owner.Login,
				Name:     *fullRepo.Parent.Name,
				URL:      *fullRepo.Parent.HTMLURL,
				BehindBy: *comparison.BehindBy,
				AheadBy:  *comparison.AheadBy,
			}
		}
		if repo.Description != nil {
			repoObject.Description = repo.Description
		}
		userRepositories = append(userRepositories, repoObject)
	}

	return userRepositories, nil
}

func (f *fetcher) GetDistinctUsers() config.UserList {
	var distinctUsers config.UserList
	for _, user := range f.Config.Users {
		if distinctUsers.Contains(user) {
			log.Warnln("Duplicate user found with Username: " + user.Name)
			continue
		}
		distinctUsers = append(distinctUsers, user)
	}

	return distinctUsers
}

func (f *fetcher) Fetch() (models.AwesomeList, error) {
	distinctUsers := f.GetDistinctUsers()
	al := models.AwesomeList{
		Users: []models.User{},
	}
	for _, user := range distinctUsers {
		log.Infoln("Fetching repos for user: " + user.Name)
		repos, err := f.GetFiveMReposForUser(user, append(f.Config.IgnoredRepos, user.IgnoredRepos...))
		if err != nil {
			return al, err
		}
		if len(repos) == 0 {
			log.Warnln("User with Username: " + user.Name + " has 0 repositories")
			continue
		}
		al.Users = append(al.Users, models.User{
			Name:         user.Name,
			Repositories: repos,
		})
	}

	return al, nil
}
