# Command package

The flamingo command package provides the *Flamingo Root Command* and allows to add additional commands under the Flamingo root command.
It is based on the popular spf13/cobra package.

## How to add new commands for the root command

Register your own commands via dingo multibindings to `*cobra.Command` inside your flamingo `module.go` file:

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

Or if you need Dingo to inject some configurations or other useful stuff then use a Dingo Provider function to bind your command:

``` 
// Configure method which belongs to dingo.Module interface.
// Responsible for module configuration (dependency injection, setting up handlers etc)
func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToProvider(MyCommand)
}



func MyCommand(router *Router, area *config.Area) *cobra.Command {
  ... 
}
   

``` 

If your module is part of a flamingo project, then you can call the command simply with:

```go run main.go myCommand```


### About the flamingo Root Command

The *Flamingo Root Command* is a `*cobra.Command` command annotated with `flamingo`.

It is normaly used by the default bootstrap of flamingo (see `flamingo/app.go`)

This is why the default outpout of a plain flamingo project (using the default app bootstrap)

Looks something like this:

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
