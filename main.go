package main

import (
	"flag"
	"fmt"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/config"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/fetcher"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/generator"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ConfigPath string
var OutputPath string

func init() {
	flag.StringVar(&ConfigPath, "config", "config.yaml", "Path to config.yaml")
	flag.StringVar(&OutputPath, "output", "awesome.md", "Path to output markdown file")
	flag.Parse()
}

func GetConfig() (*config.Config, error) {
	yFile, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(yFile))
	var conf config.Config
	err = yaml.Unmarshal(yFile, &conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

func main() {
	Config, err := GetConfig()
	if err != nil {
		panic(err)
	}

	fetch := fetcher.New(Config)
	al, err := fetch.Fetch()
	if err != nil {
		panic(err)
	}

	gen := generator.New(Config)
	err = gen.Generate(al, OutputPath)
	if err != nil {
		panic(err)
	}
}
