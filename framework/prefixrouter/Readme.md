# Prefixrouter

The prefix router overrides the "serve" command and offers additional prefixrouting.

This way you can achieve that different configuration areas are loaded based on an url prefix.

## Configuration

e.g. you can set the prefix in a configuration by using different configs for `flamingo.router.*` - for example like this:

*config/config.yml:*
```
flamingo.router.path: /en
```

You might then in another config area set it to a different prefix:
*config/de_de/config.yml:*
```
flamingo.router.path: /de
```

## Routing logic:

The prefixrouter runs in the following steps:
1. primaryHandlers: Check if (optional) available "primaryHandlers" match
2. prefixRouting: Check if the current url prefix matches one of the configured baseurls - and start routing in the matching configuration area.
3. secondaryHandlers: Check if (optional) available "secondaryHandlers" match

## Register handlers:

If you like to add your own primary or secondary handlers you can use dingo to do so:

```
injector.BindMulti((*prefixrouter.OptionalHandler)(nil)).AnnotatedWith("primaryHandlers").To(redirector{})
```

## Available Handlers

This package also provides some typical primaryHandlers that may be useful for some projects - you can activate them by configuration:

### RootRedirectHandler

A handler that redirects "/" to the configured target location.
This is useful for the prefixrouter to redirect to a default prefix.

```
prefixrouter.rootRedirectHandler.enabled: true
prefixrouter.rootRedirectHandler.redirectTarget: "/en/"
```
