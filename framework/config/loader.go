package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/ghodss/yaml"
)

// Load configuration in basedir
func Load(root *Area, basedir string) error {
	load(root, basedir, "/")

	// load additional single context file
	if os.Getenv("CONTEXTFILE") != "" {
		loadConfig(root, os.Getenv("CONTEXTFILE"))
	}

	root.GetFlatContexts()

	return nil
}

// LoadConfigFile loads a config
func LoadConfigFile(area *Area, file string) error {
	err := loadConfig(area, file)
	area.GetFlatContexts()
	return err
}

func load(area *Area, basedir, curdir string) error {
	loadConfig(area, filepath.Join(basedir, curdir, "config.yml"))
	loadRoutes(area, filepath.Join(basedir, curdir, "routes.yml"))
	if os.Getenv("CONTEXT") != "" {
		loadConfig(area, filepath.Join(basedir, curdir, "config_"+os.Getenv("CONTEXT")+".yml"))
		loadRoutes(area, filepath.Join(basedir, curdir, "routes_"+os.Getenv("CONTEXT")+".yml"))
	}
	loadConfig(area, filepath.Join(basedir, curdir, "config_local.yml"))
	loadRoutes(area, filepath.Join(basedir, curdir, "routes_local.yml"))

	for _, child := range area.Childs {
		load(child, basedir, filepath.Join(curdir, child.Name))
	}

	return nil
}

var regex = regexp.MustCompile(`%%ENV:([^%]+)%%`)

func loadConfig(area *Area, filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	config = []byte(regex.ReplaceAllStringFunc(
		string(config),
		func(a string) string { return os.Getenv(regex.FindStringSubmatch(a)[1]) },
	))

	cfg := make(Map)
	err = yaml.Unmarshal(config, &cfg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(area.Name, "loading", filename)

	if area.LoadedConfig == nil {
		area.LoadedConfig = make(Map)
	}

	area.LoadedConfig.Add(cfg)
	return nil
}

func loadRoutes(area *Area, filename string) error {
	routes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}

	err = yaml.Unmarshal(routes, &area.Routes)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println(area.Name, "loading", filename)

	return nil
}
