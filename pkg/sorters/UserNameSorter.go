package sorters

import (
	"strings"

	"github.com/iLLeniumStudios/awesome-list-generator/pkg/models"
)

type UserNameSorter []models.User

func (s UserNameSorter) Len() int {
	return len(s)
}
func (s UserNameSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s UserNameSorter) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name)
}
