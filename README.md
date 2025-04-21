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
  - [Reading Values from External Files](#reading-values-from-external-files)
  - [Options](#options)
    - [Optional Values](#optional-values)
    - [Splitting into Slices](#splitting-into-slices)
      - [Default Values](#default-values)
      - [Specifying a Custom Prefix](#specifying-a-custom-prefix)
      - [Additional Fallback Values](#additional-fallback-values)
  - [Custom Error Parsing](#custom-error-parsing)

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

## Reading Values from External Files

`minienv` supports loading variables from `.env` files by using the `WithFile(...)` option:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment

// `WithFile()` with no arguments will look for a `.env` file in the current directory
err := minienv.Load(&e, minienv.WithFile(false, "some.env", "extra.env")) 
if err != nil {
    // handle error
}
```

The first argument controls if the files are required to be there or not. `false` indicates that the load will just continue if the file / files were not found, a `true` on the other hand would raise an error if a file was not found of couldn't be parsed.

**Precedence Order:** Values from `.env`-files have a lower precedence than environment variables, therefore if a key exists in the environment and in a `.env`-file, there value in the environment takes precedence. Also, if a key exists in multiple `.env`-files, the last value takes precedence.

## Options

Minienv supports various options that can be used to control the behavior for a specfiic environment variable.

### Optional Values

By default every value is required, if no matching env variables was found or no default was specified, the load will fail with an error.
This can be changed by declaring a certain field as optional in the tag:

```go
type Environment struct {
    Port int `env:"PORT,optional"`
}

var e Environment
if err := minienv.Load(&e); err != nil {
    // handle error
}

print(e.Port) // will simply be default value of int
```

### Splitting into Slices

`minienv` supports loading data into slices and automatically splitting them on a specified token:

```go
type Environment struct {
    Ports []int `env:"PORT,split=,"` // "," is a valid split
}

var e Environment
if err := minienv.Load(&e); err != nil {
    // handle error
}
```

If no token is specified, an empty character will be used for splitting.

#### Default Values

Minienv allows you to specify default values that will be used if no value was found in the environment or through other mechanism like `WithFile()`.

```go
type Environment struct {
    Port int `env:"PORT,default=8080"`
}

var e Environment
err := minienv.Load(&e) 
if err != nil {
    // handle error
}

print(e.Port) // will be 8080 if PORT is not set otherwise
```

When setting a default for a slice, it is possible to wrap the default in `[]`-brackets, which enables the usage of `,`.

```go
type Environment struct {
    Ports []int `env:"PORT,split=,,default=[8080,8090]"`
}

var e Environment
err := minienv.Load(&e) 
if err != nil {
    // handle error
}

print(e.Ports) // [ 8080, 8090 ]
```

#### Specifying a Custom Prefix

Another option allows you to set a prefix that will be used during environment lookup. This is applied to **all** environment variables:

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

print(e.Port) // 12345
```

## Custom Error Parsing

If Minienv encounters any issues during loading, it will raise an error to the enduser. These errors are wrapped in custom error objects that allow you to react to them more precisely.

If the input to the `Load()`-function itself is invalid, Minienv will raise the predefined `ErrInvalidInput`-error:

```go
var e Environment
err := minienv.Load(e) // e is not a pointer here, therefore invalid
if err == minienv.ErrInvalidInput {
    // do something...
}
```

Additionally if Minienv fails to load a value into a certain field, for example due to a type mismatch, it will raise an error of type `LoadError`:

```go
var e Environment
err := minienv.Load(&e)
if err != nil {
    loadErr = err.(minienv.LoadError)
    // handle load error
}
```

The `LoadError` additionally exposes the affected field that failed together with the underlying error.