# Troubleshooting

## Introduction
During development one sometimes may face something, one did not expect to see. For that case we have prepared several tools, which can help you to make things working.

### Cobra Commands
Those are commands you can execute instead of `serve` to have an interesting result.

1. `routes` - dumps a list of registered routes and handler names. This may help you to understand, why your route is not triggering actions, you expect it to trigger.
1. `handler` - dumps handler names and their actions. Helps a lot in combination with `routes`.
1. `modules` - dumps an information about modules registered in flamingo and modules in child areas. This may help, when you have a big app and not sure which modules registered for which area.
1. `config` - dumps a compiled config which is currently active for your application. You can reduce dump, by adding `yml` address of a config part, you really would like to see. For example `config flamingo.session`. Heavily helps debugging if your configuration compiled from several places, for example when using helm charts.

### Flags
Flamingo utilises native go tool - `flag.FlagSet`. Flag used before `serve` slightly changes some aspects of execution, and should be used only for debug purposes.

1. `--dingo-trace-circular` - helps by giving you more information about circular injections, which of course make complete execution of your app impossible.
1. `--dingo-trace-injections` - prints what is getting injecting and which field is getting set in real time. May help you understand, what you possibly did wrong, so your application breaks during injection phase.
