package config

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
)

// LoadYaml starts to recursive read the yaml context tree
func LoadYaml(basedir string, root *Area) error {
	// load context.yml
	err := loadyaml(basedir, "/", root)
	if err != nil {
		return err
	}

	// load additional single context file
	if os.Getenv("CONTEXTFILE") != "" {
		contextfile, err := ioutil.ReadFile(os.Getenv("CONTEXTFILE"))
		if err != nil {
			return err
		}
		yaml.Unmarshal(contextfile, root)
	}

	return nil
}

func loadyaml(basedir string, curdir string, root *Area) error {
	// load context.yml
	contextfile, err := ioutil.ReadFile(path.Join(basedir, curdir, "context.yml"))
	if err == nil {
		yaml.Unmarshal(contextfile, root)
	}

	// load context_CONTEXT.yml
	contextfile, err = ioutil.ReadFile(path.Join(basedir, curdir, "context_"+os.Getenv("CONTEXT")+".yml"))
	if err == nil {
		yaml.Unmarshal(contextfile, root)
	}

	if envconfig := os.Getenv("CONFIG"); envconfig != "" {
		yaml.Unmarshal([]byte(envconfig), root)
	}

	for _, child := range root.Childs {
		err := loadyaml(basedir, path.Join(curdir, child.Name), child)
		if err != nil {
			return err
		}
	}

	return nil
}
