package context

import (
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// BUG(bastian.ike) Refactor

/// MustLoadYaml panics when LoadYaml fails
func MustLoadYaml(dir string) map[string]*Context {
	r, err := LoadYaml(dir)
	if err != nil {
		panic(err)
	}
	return r
}

// LoadYaml loads the contexts from a given folder and resolves Parent-relationships
func LoadYaml(dir string) (map[string]*Context, error) {
	basecfg, err := ioutil.ReadFile(path.Join(dir, "context.yml"))
	if err != nil {
		return nil, err
	}

	contextlist := make(map[string]string)
	yaml.Unmarshal(basecfg, &contextlist)

	result := make(map[string]*Context)

	for baseurl, name := range contextlist {
		result[name] = &Context{
			BaseUrl: baseurl,
			Name:    name,
		}

		contextConfig, err := ioutil.ReadFile(path.Join(dir, name, "config.yml"))
		if err != nil {
			return nil, err
		}

		yaml.Unmarshal(contextConfig, &result[name].Configuration)

		routing, err := ioutil.ReadFile(path.Join(dir, name, "routing.yml"))
		if err != nil {
			return nil, err
		}

		yaml.Unmarshal(routing, &result[name].Routes)
	}

	for _, context := range result {
		if context.Configuration["parent"] != "" {
			context.Parent = result[context.Configuration["parent"]]
		}
	}

	return result, nil
}
