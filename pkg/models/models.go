package models

import "strings"

type Fork struct {
	Owner    string
	Name     string
	URL      string
	AheadBy  int
	BehindBy int
}

type Repository struct {
	Name        string
	Description *string
	URL         string
	Fork        *Fork
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
