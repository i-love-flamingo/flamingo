# Troubleshooting

## Introduction
During development, one sometimes encounter something one did not expect to see. In that case, we have prepared several tools that can help you get things working.

### Cobra Commands
Those are commands you can execute instead of `serve` to have an interesting result.

1. `routes` - dumps a list of registered routes and handler names. This may help you understand, why your route is not triggering the actions you expect it to trigger.
1. `handler` - dumps handler names and their actions. It helps a lot in combination with `routes`.
1. `modules` - dumps an information about modules registered in flamingo and modules in child areas. This may help, when you have a big app and are not sure which modules are registered for which area.
1. `config` - dumps a compiled config which is currently active for your application. You can reduce the dump, by adding `yml` address of a config part, you really would like to see. For example `config flamingo.session`. This greatly helps debugging if your configuration is compiled from several places, for example when using helm charts.

### Flags
Flamingo leverages the native Go tool `flag.FlagSet` to manage command-line flags. Most flags are for debugging purposes and can be expensive on execution, especially the dingo-related flags.

1. `--dingo-trace-circular` - helps by giving you more information about circular injections, which of course makes the complete execution of your app impossible.
1. `--dingo-trace-injections` - prints what is getting injected and which field is getting set in real-time. May help you understand, what you possibly did wrong, so your application breaks during the injection phase.
1. `--flamingo-config` - helps to manually override config values, like `--flamingo.config "flamingo.session.name: foo" --flamingo.config "flamingo.otherConfig: bar"`.
1. `--flamingo-context` - sets a default execution context using a full path, like `root/mainstore/en_GB`.
1. `--flamingo-config-log` - helps to debug a configuration. During application startup, it will print all configuration compilation settings, like which configuration file is loaded in which order.
1. `--flamingo-config-cue-debug` - prints what is set in a compiled cue configuration of the project. To print everything use `.`, or to print the exact place use an address like `flamingo.session`.
1. `--dingo-inspect` - prints all the bindings that are present in a dingo injector. Can show you which values are assigned to which interface.
