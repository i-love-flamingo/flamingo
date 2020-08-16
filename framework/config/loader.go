package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"cuelang.org/go/cue/format"
	"cuelang.org/go/cue/parser"
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
		if err := root.loadConfig(false, false); err != nil {
			log.Println(err)
		}
		config.cueDebugCallback(format.Node(root.cueInstance.Lookup(config.cueDebugPath...).Syntax(), format.Simplify()))
	}
	return root.loadConfig(config.legacy, config.logLegacy)
}

func loadConfigFromBasedir(root *Area, config *LoadConfig) error {
	if err := load(root, config.basedir, "/", config); err != nil {
		return err
	}

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

func load(area *Area, basedir, curdir string, config *LoadConfig) error {
	loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config"), config.debug)
	loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config"), config.debug)
	loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes"), config.debug)
	for _, context := range strings.Split(os.Getenv("CONTEXT"), ":") {
		if context == "" {
			continue
		}
		loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config_"+context+""), config.debug)
		loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config_"+context+""), config.debug)
		loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes_"+context+""), config.debug)
	}
	loadLogged(area, loadYamlFile, filepath.Join(basedir, curdir, "config_local"), config.debug)
	loadLogged(area, loadCueFile, filepath.Join(basedir, curdir, "config_local"), config.debug)
	loadLogged(area, loadYamlRoutesFile, filepath.Join(basedir, curdir, "routes_local"), config.debug)

	for _, child := range area.Childs {
		if err := load(child, basedir, filepath.Join(curdir, child.Name), config); err != nil {
			return err
		}
	}
	return nil
}

func loadCueFile(area *Area, filename string) error {
	f, err := os.Open(filename + ".cue")
	if f != nil {
		_ = f.Close()
	}
	if err != nil {
		return nil
	}

	file, err := parser.ParseFile(filename+".cue", nil)
	if err != nil {
		return err
	}
	area.cueConfig = cueAstMergeFile(area.cueConfig, file)

	return nil
}

var regex = regexp.MustCompile(`%%ENV:([^%\n]+)%%(([^%\n]+)%%)?`)

func loadYamlFile(area *Area, filename string) error {
	config, err := ioutil.ReadFile(filename + ".yml")
	if err == nil {
		return loadYamlConfig(area, config)
	}

	config, err = ioutil.ReadFile(filename + ".yaml")
	if err == nil {
		return loadYamlConfig(area, config)
	}

	return fmt.Errorf("can not load %s.yml nor %s.yaml", filename, filename)
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
		log.Print(errorLineDebug(err, config))
		panic(err)
	}

	if area.loadedConfig == nil {
		area.loadedConfig = make(Map)
	}

	return area.loadedConfig.Add(cfg)
}

//errorLineDebug returns the lines where the error occurred (if possible)
func errorLineDebug(err error, config []byte) string {
	r, _ := regexp.Compile(": line (.*):")
	matches := r.FindStringSubmatch(err.Error())
	if len(matches) != 2 {
		return ""
	}
	line, aerr := strconv.Atoi(matches[1])
	if aerr != nil {
		return ""
	}
	errorLines := fmt.Sprintln("")
	lines := strings.Split(string(config), "\n")
	first := line - 10
	last := line + 10
	if first < 0 {
		first = 0
	}
	if last > len(lines) {
		last = len(lines)
	}
	for i := first; i < last; i++ {
		if i == line {
			errorLines = errorLines + fmt.Sprintln(">", i, lines[i])
		} else {
			errorLines = errorLines + fmt.Sprintln(" ", i, lines[i])
		}
	}
	return errorLines
}

func loadYamlRoutesFile(area *Area, filename string) error {
	routes, err := ioutil.ReadFile(filename + ".yml")
	if err == nil {
		return yaml.Unmarshal(routes, &area.Routes)
	}

	routes, err = ioutil.ReadFile(filename + ".yaml")
	if err == nil {
		return yaml.Unmarshal(routes, &area.Routes)
	}

	return fmt.Errorf("can not load %s.yml nor %s.yaml  %v", filename, filename, errorLineDebug(err, routes))
}
