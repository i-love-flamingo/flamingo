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
go run akl.go server
```

Then open: http://localhost:3210/


## Adding new dependencies

```
glide get github.com/urfave/cli
```
