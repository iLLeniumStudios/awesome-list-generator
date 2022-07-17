package config

type User struct {
	Name         string   `yaml:"name"`
	IgnoredRepos []string `yaml:"ignoredRepos"`
	SkipForks    bool     `yaml:"skipForks"`
}

type Config struct {
	Prefix          string   `yaml:"prefix"`
	MinStars        int      `yaml:"minStars"`
	MinStarsForFork int      `yaml:"minStarsForFork"`
	IgnoredRepos    []string `yaml:"ignoredRepos"`
	Users           []User   `yaml:"users"`
}
