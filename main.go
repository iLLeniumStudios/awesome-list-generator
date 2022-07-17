package main

import (
	"flag"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/config"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/fetcher"
	"github.com/iLLeniumStudios/awesome-list-generator/pkg/generator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var ConfigPath string
var OutputPath string
var Verbose bool

func init() {
	ParseFlags()
	ConfigureLogger()
}

func ParseFlags() {
	flag.BoolVar(&Verbose, "verbose", false, "Enable verbose logging")
	flag.StringVar(&ConfigPath, "config", "config.yaml", "Path to config.yaml")
	flag.StringVar(&ConfigPath, "config", "config.yaml", "Path to config.yaml")
	flag.StringVar(&OutputPath, "output", "awesome.md", "Path to output markdown file")
	flag.Parse()
}

func ConfigureLogger() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
	})

	if Verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func GetConfig() (*config.Config, error) {
	yFile, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		return nil, err
	}
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
