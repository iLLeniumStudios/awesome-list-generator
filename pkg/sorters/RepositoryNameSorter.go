package sorters

import (
	"strings"

	"github.com/iLLeniumStudios/awesome-list-generator/pkg/models"
)

type RepositoryNameSorter []models.Repository

func (s RepositoryNameSorter) Len() int {
	return len(s)
}

func (s RepositoryNameSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RepositoryNameSorter) Less(i, j int) bool {
	return strings.ToLower(s[i].Name) < strings.ToLower(s[j].Name)
}
