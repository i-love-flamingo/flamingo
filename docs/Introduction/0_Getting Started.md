# Getting started with Flamingo

## Development Setup

* 


### New? Start with [Hello World Tutorial](1_1_Tutorial Hello World.md)

---

You can use DEV context by setting the environment variable "CONTEXT"

```
export CONTEXT="dev"
```


e.g.

```
export CONTEXT="dev" && go run akl.go server
```

This will load additional config yaml files - and you can use it to point to other service urls while developing.
