# Flamingo Project Structure


```
flamingo (Project Root)
│   README.md
│   Dockerfile
│   Makefile
│   Jenkinsfile
│   glide.yaml
│   glide.lock
│
└───PROJECTNAME
│   │   PROJECTNAME.go (Main executable)
│   │
│   └───config (Main project config)
│   └───src (Project specific packages go here)       
│   
└───core (Core packages are loaded here)
│   └───auth
│   └───cms
│   └───product
│
└───framework (Framework packages go here)
│   └───router
│   └───web
│
└───om3 (OM3 related packages go here)
│
└───... (Optional additional packages can be grouped)

```

## PROJECTNAME structure

Normaly this is where all the project specific stuff belongs to.
(Name it like you want)

The Main executable should do:
* build the context tree for your project
* start the root command (delegate work to core/cmd package)

### Configuration Context

A configuration context allows to run several different flamingo powered sites in one installation.
Typical usecases are localisations and/or specific subsites.

A Context represents:
* a list of modules for each context
* configurations (that are loaded from the coresponsing "config" folder)
* a baseurl: This is the first part of the URL - and its telling flamingo which context should be loaded.
* Contexts can be nested in a tree like structure - this allows you to inherit the properties of contexts from parents to childs
