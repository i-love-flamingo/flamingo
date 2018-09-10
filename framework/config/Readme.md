# Config module

This module provides features to read and merge flamingo configurations.
Also it provides templatefunctions to access configurations in your template.

Configurations are defined and used by the individual modules. 
The modules should come with a documentation which configurations/featureflags they support.


## Basics
Configurations are yml files located in **config** folder.

The configuration syntax is to specify areas either with `.` or as yaml maps:

```yaml
foo:
  bar: x
```

is the same as

```yaml
foo.bar: x
```

Configuration can be used:

* either by the `config()` templatefunction in your template
* or via dependency injection from dingo

## Loaded Configuration files
The following configuration files will be loaded from the given `config` folder:

* config.yml
* routes.yml
* optional: config_($CONTEXT).yml
* optional: routes_($CONTEXT).yml
* config_local.yml
* routes_local.yml

You can set different CONTEXT with the environment variable *CONTEXT* and this will cause flamingo to load additional configuration files.

e.g. starting flamingo with
```bash
export CONTEXT="dev" && go run project.go serve
```
Will cause flamingo to additionaly load the configfile "config/config_dev.yml"


Configuration values can also be read from environment variables with the syntax:

```yaml
auth.secret: '%%ENV:KEYCLOAK_SECRET%%'
```


### Injecting configurations
Asking for either a concrete value via e.g. `foo.bar` is possible, as well as getting a whole `config.Map` instance by a partially-selector, e.g. `foo`.
This would be a Map with element `bar`.

For example:
```go
Example struct {
  Title           string      `inject:"config:mymodul.title"
  CompleteConfig  config.Map  `inject:"config:mymodul"
  Amount          float64     `inject:"config:mymodul.amount"
  Flag            bool        `inject:"config:mymodul.flag"
}
```

Deeply nested config maps can be marshaled into structs for convenience.

The result struct must match exactly the structure and types of the config map and all fields must be exported.

```go
err := m.MarshalTo(&result)
```
