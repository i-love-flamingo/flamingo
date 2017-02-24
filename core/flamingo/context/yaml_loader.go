package context

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// LoadYaml starts to recursive read the yaml context tree
func LoadYaml(basedir string, root *Context) error {
	// load context.yml
	return loadyaml(basedir, "/", root)
}

func loadyaml(basedir string, curdir string, root *Context) error {
	// load context.yml
	contextfile, _ := ioutil.ReadFile(path.Join(basedir, curdir, "context.yml"))

	yaml.Unmarshal(contextfile, root)

	for _, child := range root.Childs {
		err := loadyaml(basedir, path.Join(curdir, child.Name), child)
		if err != nil {
			return err
		}
	}

	return nil
}
