package fetcher

import (
	"context"

	"github.com/gammazero/workerpool"

	"os"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/config"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/models"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/utils"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type FilterType int64

const (
	Inclusion FilterType = iota
	Exclusion
)

type RepoFilter struct {
	Repos utils.StringList
	Type  FilterType
}

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

func (f *fetcher) GetFiveMReposForUser(user config.User, filter RepoFilter) ([]models.Repository, error) {
	var userRepositories []models.Repository
	repos, err := f.GetGithubReposForUsername(user.Name)
	if err != nil {
		return nil, err
	}
	for _, repo := range repos {
		if *repo.Private || *repo.StargazersCount < f.Config.MinStars || *repo.Disabled || *repo.Name == user.Name {
			continue
		}
		switch filter.Type {
		case Exclusion:
			if filter.Repos.Contains(*repo.Name) {
				continue
			}
		case Inclusion:
			if !filter.Repos.Contains(*repo.Name) {
				continue
			}
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
	wp := workerpool.New(f.Config.NumWorkers)
	count := 0
	distinctUsers := f.GetDistinctUsers()
	returnedUsers := make(chan models.User, len(distinctUsers))
	al := models.AwesomeList{
		Users: []models.User{},
	}

	for _, user := range distinctUsers {
		user := user
		wp.Submit(func() {
			log.Infoln("Fetching repos for user: " + user.Name)
			var filter RepoFilter
			if len(user.IncludeRepos) != 0 {
				filter.Type = Inclusion
				filter.Repos = user.IncludeRepos
			} else {
				filter.Type = Exclusion
				filter.Repos = append(f.Config.ExcludeRepos, user.ExcludeRepos...)
			}
			repos, err := f.GetFiveMReposForUser(user, filter)
			if err != nil {
				panic(err)
			}
			if len(repos) == 0 {
				log.Warnln("User with Username: " + user.Name + " has 0 repositories")
			} else {
				count += 1
				returnedUsers <- models.User{
					Name:         user.Name,
					Repositories: repos,
				}
			}
		})
	}

	wp.StopWait()

	for i := 0; i < count; i++ {
		al.Users = append(al.Users, <-returnedUsers)
	}

	return al, nil
}
