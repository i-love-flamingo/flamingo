package config

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"
)

var debugLog bool

func init() {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.SetOutput(new(bytes.Buffer))
	fs.BoolVar(&debugLog, "flamingo-config-log", false, "enable flamingo config loader logging")
	if err := fs.Parse(os.Args[1:]); err == flag.ErrHelp {
		fs.SetOutput(os.Stderr)
		fs.PrintDefaults()
	}
}

// Load configuration in basedir
func Load(root *Area, basedir string) error {
	load(root, basedir, "/")

	// load additional single context file
	for _, file := range strings.Split(os.Getenv("CONTEXTFILE"), ":") {
		if file == "" {
			continue
		}
		loadConfig(root, file)
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
	for _, context := range strings.Split(os.Getenv("CONTEXT"), ":") {
		if context == "" {
			continue
		}
		loadConfig(area, filepath.Join(basedir, curdir, "config_"+context+".yml"))
		loadRoutes(area, filepath.Join(basedir, curdir, "routes_"+context+".yml"))
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
		if debugLog {
			log.Println(err)
		}
		return err
	}

	config = []byte(regex.ReplaceAllFunc(
		config,
		func(a []byte) []byte { return []byte(os.Getenv(string(regex.FindSubmatch(a)[1]))) },
	))

	cfg := make(Map)
	err = yaml.Unmarshal(config, &cfg)
	if err != nil {
		if debugLog {
			log.Println(err)
		}
		return err
	}

	if debugLog {
		log.Println(area.Name, "loading", filename)
	}

	if area.LoadedConfig == nil {
		area.LoadedConfig = make(Map)
	}

	area.LoadedConfig.Add(cfg)
	return nil
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
		if debugLog {
			log.Println(err)
		}
		return err
	}

	if debugLog {
		log.Println(area.Name, "loading", filename)
	}

	return nil
}
