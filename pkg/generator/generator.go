package generator

import (
	"encoding/json"
	"errors"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/config"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/models"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type Generator interface {
	Generate(al models.AwesomeList, outputPath string, outputFormat string) error
}

type generator struct {
	Config *config.Config
}

func New(conf *config.Config) Generator {
	return &generator{
		Config: conf,
	}
}

func (g *generator) generateJSON(al models.AwesomeList, outputPath string) error  {
	var allRepos []models.Repository
	for _, user := range al.Users {
		allRepos = append(allRepos, user.Repositories...)
	}
	file, _ := json.MarshalIndent(allRepos, "", "\t")
	return ioutil.WriteFile(outputPath, file, 0644)
}

func (g *generator) generateMarkdown(al models.AwesomeList, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	sort.Sort(al.Users)

	f.WriteString(g.Config.Prefix)

	defer f.Close()

	f.WriteString("\n")

	for _, user := range al.Users {
		f.WriteString("### " + user.Name + " <a name=\"" + user.Name + "\"></a>\n")
		for _, repo := range user.Repositories {
			f.WriteString("- [" + repo.Name + "](" + repo.URL + ")")
			if repo.Fork != nil {
				f.WriteString(" ([Original](" + repo.Fork.URL + ") :green_circle: +" + strconv.Itoa(repo.Fork.AheadBy) + " :red_circle: -" + strconv.Itoa(repo.Fork.BehindBy) + "</span>)")
			}
			if repo.Description != nil {
				f.WriteString(" - " + *repo.Description)
			}
			f.WriteString(" ![GitHub stars](https://img.shields.io/github/stars/" + user.Name + "/" + repo.Name + ".svg?style=social&label=Stars&maxAge=2592000)\n")
		}
		f.WriteString("\n")
	}

	return nil
}

func (g *generator) Generate(al models.AwesomeList, outputPath string, outputFormat string) error {

	switch outputFormat {
	case "markdown":
		return g.generateMarkdown(al, outputPath)
	case "json":
		return g.generateJSON(al, outputPath)
	default:
		return errors.New("invalid output format")
	}
}
