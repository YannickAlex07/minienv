# Minimal Reflection-Based Environment Management for Go

[![codecov](https://codecov.io/gh/YannickAlex07/minienv/branch/main/graph/badge.svg?token=VHXLuQARRp)](https://codecov.io/gh/YannickAlex07/minienv)
[![Go Reference](https://pkg.go.dev/badge/github.com/yannickalex07/minienv.svg)](https://pkg.go.dev/github.com/yannickalex07/minienv)

`minienv` is a minimal libary that makes it easy to work with environment variables in Go. It is heavily inspired by [`netflix/go-env`](https://github.com/Netflix/go-env) and Pythons [`pydantic/BaseSettings`](https://docs.pydantic.dev/latest/concepts/pydantic_settings/) and combines reading from `.env` files and reflection based parsing of environment variables.

Add it with the following command:

```
go get github.com/yannickalex07/minienv
```
- [Minimal Reflection-Based Environment Management for Go](#minimal-reflection-based-environment-management-for-go)
  - [Getting Started](#getting-started)
      - [Optional Values](#optional-values)
      - [Default Values](#default-values)
      - [Reading `.env`-Files](#reading-env-files)
      - [Additional Fallback Values](#additional-fallback-values)
      - [Specifying a Custom Prefix](#specifying-a-custom-prefix)

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

#### Optional Values

By default every value is required, so if no matching env variables was found or no default is specified, the load will fail with an error.
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

#### Default Values

Minienv allows you to specify default values that will be used if no value was found in the environment or specified through a fallback like `WithFile()` or `WithFallbackValues()`.

```go
type Environment struct {
    Port int `env:"PORT,default=8080"`
}

var e Environment

// `WithFile()` with no arguments will look for a `.env` file in the current directory
err := minienv.Load(&e) 
if err != nil {
    // handle error
}

print(e.Port) // will be 8080 if PORT is not set in the environment
```

#### Reading `.env`-Files

`minienv` additionally supports loading variables from `.env` files by using the `WithFile(...)` option:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment

// `WithFile()` with no arguments will look for a `.env` file in the current directory
err := minienv.Load(&e, minienv.WithFile(false)) 
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

err := minienv.Load(&e, minienv.WithFile(true, "database.env", "extra.env"))
if err != nil {
    // handle error
}
```

The first argument controls if the files are required to be there or not. `false` indicates that the load will just continue if the file / files were not found, a `true` on the other hand would raise an error if a file was not found of couldn't be parsed.

**Precedence Order:** Values from `.env`-files have a lower precedence than environment variables, therefore if a key exists in the environment and in a `.env`-file, there value in the environment takes precedence. Also, if a key exists in multiple `.env`-files, the last value takes precedence.

#### Additional Fallback Values

Another option that `minienv` provides is to supply custom fallback values that might be sourced from somewhere completely else:

```go
type Environment struct {
    Port int `env:"PORT"`
}

values := map[string]string{
    "PORT": "12345"
}

var e Environment
err := minienv.Load(&e, minienv.WithFallbackValues(values))
if err != nil {
    // handle error
}
```

#### Specifying a Custom Prefix

Another option allows you to set a prefix that will be used during environment lookup:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment
err := minienv.Load(&e, minienv.WithPrefix("APP_")) // will cause a lookup for APP_PORT
if err != nil {
    // handle error
}
```

This prefix is also applied to keys from `.env`-files as well as additional fallback values, however only if the key does not already contain the prefix.