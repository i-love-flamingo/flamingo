package config

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/format"
	"github.com/ghodss/yaml"
)

type (
	// LoadConfig provides configuration for the loader
	LoadConfig struct {
		legacy           bool
		logLegacy        bool
		additionalConfig []string
		basedir          string
		debug            bool
		cueDebugPath     []string
		cueDebugCallback func([]byte, error)
	}

	// LoadOption to be passed to Load(, ...)
	LoadOption func(*LoadConfig)
)

// DebugLog enables/disabled detailed debug logging
func DebugLog(debug bool) LoadOption {
	return func(config *LoadConfig) {
		config.debug = debug
	}
}

// CueDebug enables a cue.Instance debugger. This is part of a dev-api and might change!
func CueDebug(path []string, callback func([]byte, error)) LoadOption {
	return func(config *LoadConfig) {
		config.cueDebugPath = path
		config.cueDebugCallback = callback
	}
}

// LegacyMapping controls if flamingo legacy config mapping happens
func LegacyMapping(mapLegacy, logLegacy bool) LoadOption {
	return func(config *LoadConfig) {
		config.legacy = mapLegacy
		config.logLegacy = logLegacy
	}
}

// AdditionalConfig adds additional config values (yaml strings) to the config
func AdditionalConfig(addtionalConfig []string) LoadOption {
	return func(config *LoadConfig) {
		config.additionalConfig = append(config.additionalConfig, addtionalConfig...)
	}
}

// Load configuration in basedir
func Load(root *Area, basedir string, options ...LoadOption) error {
	config := &LoadConfig{
		legacy:  true,
		basedir: basedir,
	}
	for _, option := range options {
		option(config)
	}
	if err := loadConfigFromBasedir(root, config); err != nil {
		return err
	}
	if config.cueDebugCallback != nil {
		_ = root.loadConfig(false, false)
		config.cueDebugCallback(format.Node(root.cueInstance.Lookup(config.cueDebugPath...).Syntax(), format.Simplify()))
	}
	return root.loadConfig(config.legacy, config.logLegacy)
}

func loadConfigFromBasedir(root *Area, config *LoadConfig) error {
	load(root, config.basedir, "/", config.debug)

	// load additional single context file
	for _, file := range strings.Split(os.Getenv("CONTEXTFILE"), ":") {
		file = strings.TrimSuffix(file, filepath.Ext(file))
		if file == "" {
			continue
		}
		loadLogged(root, loadYamlFile, file, config.debug)
		loadLogged(root, loadCueFile, file, config.debug)
	}

	for _, add := range config.additionalConfig {
		if config.debug {
			log.Printf("Loading %q", add)
		}
		if err := loadYamlConfig(root, []byte(add)); err != nil {
			return err
		}
	}

	return nil
}

// LoadConfigFile loads a config
// Deprecated: do not arbitrarily load anything anymore, use Area.Load
func LoadConfigFile(area *Area, file string) error {
	log.Println("WARNING! config.LoadConfigFile is deprecated!")

	if err := loadYamlFile(area, file); err != nil {
		return err
	}
	if err := loadCueFile(area, file); err != nil {
		return err
	}
	return nil
}

func loadLogged(area *Area, loader func(*Area, string) error, filename string, debug bool) {
	if debug {
		log.Printf("Loading %q", filename)
	}
	if err := loader(area, filename); err != nil && debug {
		log.Printf("Error: %s", err)
	}
}

func load(area *Area, basedir, curdir string, debug bool) {
	loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config"), debug)
	loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config"), debug)
	loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes"), debug)
	for _, context := range strings.Split(os.Getenv("CONTEXT"), ":") {
		if context == "" {
			continue
		}
		loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config_"+context+""), debug)
		loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config_"+context+""), debug)
		loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes_"+context+""), debug)
	}
	loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config_local"), debug)
	loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config_local"), debug)
	loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes_local"), debug)

	for _, child := range area.Childs {
		load(child, basedir, filepath.Join(curdir, child.Name), debug)
	}
}

func loadCueFile(area *Area, filename string) error {
	f, err := os.Open(filename + ".cue")
	if f != nil {
		_ = f.Close()
	}
	if err != nil {
		return nil
	}

	if area.cueBuildInstance == nil {
		cueContext := build.NewContext()
		area.cueBuildInstance = cueContext.NewInstance(area.Name, nil)
	}

	return area.cueBuildInstance.AddFile(filename+".cue", nil)
}

var regex = regexp.MustCompile(`%%ENV:([^%\n]+)%%(([^%\n]+)%%)?`)

func loadYamlFile(area *Area, filename string) error {
	config, err := ioutil.ReadFile(filename + ".yml")
	if err != nil {
		return err
	}
	return loadYamlConfig(area, config)
}

func loadYamlConfig(area *Area, config []byte) error {
	config = regex.ReplaceAllFunc(
		config,
		func(a []byte) []byte {
			value := os.Getenv(string(regex.FindSubmatch(a)[1]))
			if value == "" {
				value = string(regex.FindSubmatch(a)[3])
			}
			return []byte(value)
		},
	)

	cfg := make(Map)
	if err := yaml.Unmarshal(config, &cfg); err != nil {
		panic(err)
	}

	if area.loadedConfig == nil {
		area.loadedConfig = make(Map)
	}

	return area.loadedConfig.Add(cfg)
}

func loadYamlRoutesFile(area *Area, filename string) error {
	routes, err := ioutil.ReadFile(filename + ".yml")
	if err != nil {
		return err
	}
	return yaml.Unmarshal(routes, &area.Routes)
}
