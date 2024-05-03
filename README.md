# Minimal Environment Management for Go

[![codecov](https://codecov.io/gh/YannickAlex07/minienv/branch/main/graph/badge.svg?token=VHXLuQARRp)](https://codecov.io/gh/YannickAlex07/minienv)
[![Go Reference](https://pkg.go.dev/badge/github.com/yannickalex07/minienv.svg)](https://pkg.go.dev/github.com/yannickalex07/minienv)

`minienv` is a minimal libary to work with environment variables. It is heavily inspired by `netflix/go-env` and Pythons `pydantic/BaseSettings` and combines reading from `.env` files and reflection based parsing of environment variables.

Add it with the following command:

```
go get github.com/yannickalex07/minienv
```

## Getting Started

Using `minienv` is quite simple, just create a struct and annotate it with `env:""` tags:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment
if err := minienv.Load(&e); err != nil {
    // handle error
}

print(e.Port) // will equal to whatever the PORT env variable is set to
```

### Optional Values

By default every value is required, so if no matching env variables was found, the load will fail with an error.
This can be changed by declaring a certain field as optional in the tag:

```go
type Environment struct {
    Port int `env:"PORT,optional"`
}

var e Environment
if err := minienv.Load(&e); err != nil {
    // handle error
}

print(e.Port) // will be the default value of PORT was not set
```


### Reading `.env`-Files

`minienv` additionally supports loading variables from `.env` files by using the `WithFile(...)` and `WithRequiredFile(...)` options:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment

// `WithFile()` with no arguments will look for a `.env` file in the current directory
err := minienv.Load(&e, minienv.WithFile()) 
if err != nil {
    // handle error
}
```

Alternatively you can specify one or multiple explicit files:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment

err := minienv.Load(&e, minienv.WithFile("database.env", "extra.env"))
if err != nil {
    // handle error
}
```

The difference between `WithFile()` and `WithRequiredFile()` is simply that `WithRequiredFile()` will fail if any error occurs, raising the error up. `WithFile()` will essentially just silently fail and will not raise any error.

### Additional Overrides

Another option that `minienv` provides is to supply custom overrides that might be sourced from somewhere completely else:

```go
type Environment struct {
    Port int `env:"PORT"`
}

overrides := map[string]string{
    "PORT": "12345"
}

var e Environment
err := minienv.Load(&e, minienv.WithOverrides(overrides))
if err != nil {
    // handle error
}
```