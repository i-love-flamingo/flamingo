package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

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

func load(area *Area, basedir, curdir string) error {
	loadConfig(area, path.Join(basedir, curdir, "config.yml"))
	loadRoutes(area, path.Join(basedir, curdir, "routes.yml"))
	loadConfig(area, path.Join(basedir, curdir, "config_"+os.Getenv("CONTEXT")+".yml"))
	loadRoutes(area, path.Join(basedir, curdir, "routes_"+os.Getenv("CONTEXT")+".yml"))
	loadConfig(area, path.Join(basedir, curdir, "config_local.yml"))
	loadRoutes(area, path.Join(basedir, curdir, "routes_local.yml"))

	for _, child := range area.Childs {
		load(child, basedir, path.Join(curdir, child.Name))
	}

	return nil
}

func loadConfig(area *Area, filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err)
		return err
	}

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
