# Command package

The Flamingo command package provides the *Flamingo root command* and allows to add additional commands under the Flamingo root command.
It is based on the popular [spf13/cobra](https://github.com/spf13/cobra) package.

## How to add new commands for the root command

Register your own commands via Dingo multibindings to `*cobra.Command` inside your Flamingo `module.go` file:

E.g.:
```go
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToInstance(myCommand())
}

func myCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "myCommand",
		Short: "myCommand short desc",
		Run: func(cmd *cobra.Command, args []string) {
      doSomething()
		},
	}
	return cmd
}

```

Or, if you need Dingo to inject some configurations or other useful stuff then use a Dingo provider function to bind your command:

```go
// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToProvider(MyCommand)
}

// MyCommand gets called by Dingo and all arguments are resolved and injected
func MyCommand(router *Router, area *config.Area) *cobra.Command {
  ... 
}
``` 

If your module is part of a Flamingo project, then you can call the command simply with:

```bash
go run main.go myCommand
```

### About the Flamingo root command

The *Flamingo root command* is a `*cobra.Command` command annotated with `flamingo`.

It is normally used by the default bootstrap of Flamingo (see `flamingo/app.go`)

This is why the default output of a plain Flamingo project (using the default app bootstrap) looks like this:

```sh
$ go run main.go

Flamingo main

Usage:
  main [command]

Available Commands:
  config      Config dump
  help        Help about any command
  serve       Default serve command - starts on Port 3322

Flags:
  -h, --help   help for main

Use "main [command] --help" for more information about a command.
```

## Adding persistent flags to the root command

You can add persistent [flags](https://github.com/spf13/cobra#flags) to Flamingo's root command by multi-binding a 
[`FlagSet`](https://godoc.org/github.com/spf13/pflag#FlagSet) to `*pflag.FlagSet`:

```go
// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti((*pflag.FlagSet)(nil)).ToInstance(someFlagSet)
}

```


