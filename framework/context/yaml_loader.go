package context

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
)

// LoadYaml starts to recursive read the yaml context tree
func LoadYaml(basedir string, root *Context) error {
	// load context.yml
	return loadyaml(basedir, "/", root)
}

func loadyaml(basedir string, curdir string, root *Context) error {
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

	for _, child := range root.Childs {
		err := loadyaml(basedir, path.Join(curdir, child.Name), child)
		if err != nil {
			return err
		}
	}

	return nil
}
