# Flamingo cmd package

Based on spf13/cobra.

## Usage

Register your own commands via dingo multibindings to `*cobra.Command`.

Start the root command by injecting `*cobra.Command` annotated with `flamingo`.

Includes the RootCommand and some helpful Commands.

Just run your main project file to see the list of registered commands.
