# Flamingo

## Docs

Read the extensive documentation that is part of Flamingo.

After checkout you need to type:

`make doc`
This will start building and previeweing the mkdocs based documentation in a Docker container.

To view the docs open  [Docs](http://localhost:8000)

(To end the Docker container again use `docker kill`)

-----------------

# Depricated

See Docs for most recent documentation

## Prepare / Requirements

https://github.com/Masterminds/glide is used for dependencies.

```
brew install glide
```

## Fetch dependencies

```
glide install
```

## Run

Compile frontend:

```
cd akl/frontend
./build.sh
```

```
CONTEXT=dev go run akl.go serve
```

Then open: [localhost:3210](http://localhost:3210/)


## Adding new dependencies

```
glide get github.com/urfave/cli
```
