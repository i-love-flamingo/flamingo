# Config module

This module provides features to read and merge Flamingo configurations.
Also it provides template functions to access configurations in your template.

Configurations are defined and used by the individual modules. 
The modules should come with a documentation which configurations/feature-flags they support.

## Basics
Configurations are yml files located in `config` folder.

The configuration syntax is to specify sections either with `.` or as yaml maps:

```yaml
foo:
  bar: x
```

is the same as

```yaml
foo.bar: x
```

Configuration values can also be read from environment variables during the loading process with the syntax:

```yaml
auth.secret: '%%ENV:KEYCLOAK_SECRET%%'
```

or 

```yaml
auth.secret: '%%ENV:KEYCLOAK_SECRET%%default_value%%'
```

In the second case, Flamingo falls back to `default_value` if the environment variable is not set or empty.


Configuration can be used:

* either by the `config()` templatefunction in your template
* or via dependency injection from Dingo

## Loaded Configuration files
The following configuration files will be loaded from `config` folder:

* config.yml
* routes.yml
* optional: config_($CONTEXT).yml
* optional: routes_($CONTEXT).yml
* config_local.yml
* routes_local.yml

You can set different contexts with the environment variable `CONTEXT` and this will cause Flamingo to load additional configuration files.

e.g. starting Flamingo with
```bash
CONTEXT="dev" go run project.go serve
```
Will cause Flamingo to additionally load the config file `config/config_dev.yml`

You can also load multiple extra configuration files - e.g. starting Flamingo with
```bash
CONTEXT="dev:testdata" go run project.go serve
```
Will cause Flamingo to additionally load the config files "config/config_dev.yml" and "config/config_testdata.yml" in the given order.

### Additional configuration files from outside

Flamingo can load multiple additional yaml/cue files, which must be given in the environment variable `CONTEXTFILE`, separated by `:`.

The files can be given by using relative paths from the working directory or absolute paths.

```bash
CONTEXTFILE="../../myCfg.yml:/var/flamingo/cfg/main.yml:/var/flamingo/cfg/additional.cue" go run project.go serve
```

### Additional temporary configuration

You can set any configuration value within the run command by using the `--flamingo-config` flag:

```bash
go run project.go serve --flamingo-config "auth.secret: mySecret" --flamingo-config "other.secret: mySecret"
```

The flag's values are expected to be valid yaml. The flag can be used multiple times.

Configurations provided via `--flamingo-config` flag overwrite all values provided in yaml files.

### Priority of configuration

If multiple sources define the same configuration key, the value from the last loaded source is taken.
The order of loading is:

1. All files from `config` directory
  1. config.yml
  1. routes.yml
  1. All context files given in the environment variable `CONTEXT`
    1. config_($CONTEXT).yml
    1. routes_($CONTEXT).yml
  1. config_local.yml
  1. routes_local.yml
1. All files given in the environment variable `CONTEXTFILE`
1. All values given via `--flamingo-config` flag

### Debugging configuration loading

By stating `--flamingo-config-log`, you can enable the configuration loader's debug log, which prints all handled files 
to the output using go's `log` package, because the `flamingo.Logger` is not available yet in this early state of bootstrapping.


### Injecting configurations
Asking for either a concrete value via e.g. `foo.bar` is possible, as well as getting a whole `config.Map` instance by a partially-selector, e.g. `foo`.
This would be a Map with element `bar`.

All configuration values are registered as Dingo annotated binding and can be requested using the `inject` tag.
To get this in the arguments of the `Inject` function or a dingo provider, you will have to wrap it by an anonymous struct. 

For example:
```go
// Inject dependencies
func (m *Module) Inject(
	cfg *struct {
	CompleteConfig config.Map `inject:"config:mymodule"`
	Title          string     `inject:"config:mymodule.title"`
	Amount         int        `inject:"config:mymodule.amount"`
	Flag           bool       `inject:"config:mymodule.flag"`
},
) *Module {
	if cfg != nil {
		m.title = cfg.Title
		m.amount = cfg.Amount
		m.flag = cfg.Flag
		m.cfg = cfg.CompleteConfig
	}

	return m
}
```

Deeply nested config maps can be marshaled into structs for convenience.

The result struct must match exactly the structure and types of the config map and all fields must be exported.

```go
err := m.MarshalTo(&result)
```

## Using multiple configuration areas:
A Flamingo application can have multiple `config.Area` - that is essentially useful for localisation.
See [Flamingo Bootstrap](../1. Flamingo Basics/7. Flamingo Bootstrap.md)

# Convert Yaml to Cue

```
sed "s/\'%%ENV:\(.*\)%%\(.*\)%%\'/*flamingo.os.env.\1 | \"\2\"/g"
sed "s/\"%%ENV:\(.*\)%%\(.*\)%%\"/*flamingo.os.env.\1 | \"\2\"/g"
sed "s/\'%%ENV:\(.*\)%%\'/flamingo.os.env.\1/g"
sed "s/\"%%ENV:\(.*\)%%\"/flamingo.os.env.\1/g"
```
