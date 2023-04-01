package models

import (
	"strings"

	"github.com/google/go-github/v45/github"
)

type Fork struct {
	Owner    string
	Name     string
	URL      string
	AheadBy  int
	BehindBy int
}

type Repository struct {
	Owner       string
	Name        string
	Description *string `json:",omitempty"`
	URL         string
	Fork        *Fork    `json:",omitempty"`
	Tags        []string `json:",omitempty"`
	Stars       int
	LastUpdated github.Timestamp
	LastPushed  github.Timestamp
	Forks       int
	OpenIssues  int
	Archived    bool
}

type User struct {
	Name         string
	Repositories []Repository
}

type AwesomeList struct {
	Users userList
}
type userList []User

func (u userList) Len() int {
	return len(u)
}

func (u userList) Less(i, j int) bool {
	return strings.ToLower(u[i].Name) < strings.ToLower(u[j].Name)
}

func (u userList) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}
