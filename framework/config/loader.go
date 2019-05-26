package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
)

var (
	// DebugLog flag
	debugLog bool
	// AdditionalConfig to be loaded
	additionalConfig []string
	once             = sync.Once{}
)

// Load configuration in basedir
func Load(root *Area, basedir string) error {
	once.Do(func() {
		// we have to parse the flagset because config loading takes place before rootCmd's execute
		flagSet.Parse(os.Args[1:])
	})

	load(root, basedir, "/")

	// load additional single context file
	for _, file := range strings.Split(os.Getenv("CONTEXTFILE"), ":") {
		if file == "" {
			continue
		}
		if err := loadConfigFile(root, file); err != nil {
			return err
		}
	}

	for _, add := range additionalConfig {
		if debugLog {
			log.Printf("Loading %q", add)
		}
		if err := loadConfig(root, []byte(add)); err != nil {
			return err
		}
	}

	_, err := root.GetFlatContexts()
	return err
}

// LoadConfigFile loads a config
func LoadConfigFile(area *Area, file string) error {
	if err := loadConfigFile(area, file); err != nil {
		return err
	}
	_, err := area.GetFlatContexts()
	return err
}

func load(area *Area, basedir, curdir string) {
	loadConfigFile(area, filepath.Join(basedir, curdir, "config.yml"))
	loadRoutes(area, filepath.Join(basedir, curdir, "routes.yml"))
	for _, context := range strings.Split(os.Getenv("CONTEXT"), ":") {
		if context == "" {
			continue
		}
		loadConfigFile(area, filepath.Join(basedir, curdir, "config_"+context+".yml"))
		loadRoutes(area, filepath.Join(basedir, curdir, "routes_"+context+".yml"))
	}
	loadConfigFile(area, filepath.Join(basedir, curdir, "config_local.yml"))
	loadRoutes(area, filepath.Join(basedir, curdir, "routes_local.yml"))

	for _, child := range area.Childs {
		load(child, basedir, filepath.Join(curdir, child.Name))
	}
}

var regex = regexp.MustCompile(`%%ENV:([^%\n]+)%%(([^%\n]+)%%)?`)

func loadConfigFile(area *Area, filename string) error {
	config, err := ioutil.ReadFile(filename)
	if err != nil {
		if debugLog {
			log.Println(err)
		}
		return err
	}
	if debugLog {
		log.Println(area.Name, "loading", filename)
	}
	return loadConfig(area, config)
}

func loadConfig(area *Area, config []byte) error {
	config = []byte(regex.ReplaceAllFunc(
		config,
		func(a []byte) []byte {
			value := os.Getenv(string(regex.FindSubmatch(a)[1]))
			if value == "" {
				value = string(regex.FindSubmatch(a)[3])
			}
			return []byte(value)
		},
	))

	cfg := make(Map)
	err := yaml.Unmarshal(config, &cfg)
	if err != nil {
		panic(err)
	}

	if area.LoadedConfig == nil {
		area.LoadedConfig = make(Map)
	}

	return area.LoadedConfig.Add(cfg)
}

func loadRoutes(area *Area, filename string) error {
	routes, err := ioutil.ReadFile(filename)
	if err != nil {
		if debugLog {
			log.Println(err)
		}
		return err
	}

	err = yaml.Unmarshal(routes, &area.Routes)
	if err != nil {
		panic(err)
	}

	if debugLog {
		log.Println(area.Name, "loading", filename)
	}

	return nil
}
