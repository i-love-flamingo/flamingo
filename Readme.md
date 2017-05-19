# Flamingo

## Docs

`docker run --rm -v $(pwd):/work -p 8000:8000 python bash -c 'pip install mkdocs; cd /work; mkdocs serve --dev-addr=0.0.0.0:8000'`

Then view the docs at [Docs](http://localhost:8000)

# Development

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
