package config

type User struct {
	Name         string   `yaml:"name"`
	ExcludeRepos []string `yaml:"excludeRepos"`
	IncludeRepos []string `yaml:"includeRepos"`
	SkipForks    bool     `yaml:"skipForks"`
}

type Config struct {
	Prefix          string   `yaml:"prefix"`
	NumWorkers      int      `yaml:"numWorkers"`
	MinStars        int      `yaml:"minStars"`
	MinStarsForFork int      `yaml:"minStarsForFork"`
	ExcludeRepos    []string `yaml:"excludeRepos"`
	Users           UserList `yaml:"users"`
}

type UserList []User

func (ul UserList) Contains(u User) bool {
	for _, user := range ul {
		if user.Name == u.Name {
			return true
		}
	}
	return false
}
