# Flamingo Project Structure


```
flamingo (Project Root)
│   README.md
│   Dockerfile
│   Makefile
│   Jenkinsfile
│
└───PROJECTNAME
│   │   glide.*
│   │   PROJECTNAME.go (Main executable)
│   │
│   └───config (Main project config)
│   └───src (Project specific packages go here)
│       
│   
└───core (Core packages are loaded here)
│   │   glide.*
│   └───cms
└───framework (Framework packages go here)
└───... (Optional additional packages can be grouped)

```

## PROJECTNAME structure

Normaly this is where all the project specific stuff belongs to.
(Name it like you want)

The Main executable should do:
* build the context tree for your project
* start the root command (delegate work to core/cmd package)

### Context
A Context allows to run several different flamingo powered sites in one installation.
Typical usecases are localisations and/or specific subsites.

A Context represents:
* a list of packages registerfuncs
* configurations (that are loaded from the coresponsing "config" folder)
* a baseurl: This is the first part of the URL - and its telling flamingo which context should be loaded.
* Contexts can be nested in a tree like structure - this allows you to inherit the properties of contexts from parents to childs

Example init func in project exectutable:

```go
func init() {
	context.RootContext = context.New(
		"root",
		[]di.RegisterFunc{
			corecms.Register,
			pug_template.Register,
			framework.Register,
			base.Register,
			brand.Register,
			coreproduct.Register,
			internalmock.Register,
			//aklmock.Register,
			profiler.Register,
			requestlogger.Register,
			RegisterAKL,
		},
		"",
		context.New(
			"mainstore", nil,
			"",
			context.New("de", nil, "/de",),
			context.New("en", nil, "/en",),
		),
		context.New(
			"lounges", nil,
			"",
			context.New(
				"lh", nil,
				"",
				context.New("de", nil,"loungehost/lounge/de",),
				context.New("en", nil,"loungehost/lounge/en",),
			),
		),
	)
	context.LoadYaml("config", context.RootContext)
}
```
