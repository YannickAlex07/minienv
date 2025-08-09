# Minimal Reflection-Based Environment Management for Go

[![codecov](https://codecov.io/gh/YannickAlex07/minienv/branch/main/graph/badge.svg?token=VHXLuQARRp)](https://codecov.io/gh/YannickAlex07/minienv)
[![Go Reference](https://pkg.go.dev/badge/github.com/yannickalex07/minienv.svg)](https://pkg.go.dev/github.com/yannickalex07/minienv)

`minienv` is a minimal libary that makes it easy to work with environment variables in Go. It is heavily inspired by [`netflix/go-env`](https://github.com/Netflix/go-env) and Pythons [`pydantic/BaseSettings`](https://docs.pydantic.dev/latest/concepts/pydantic_settings/) and combines reflection based parsing of environment variables with reading from `.env` files.

Add it with the following command:

```
go get github.com/yannickalex07/minienv
```
- [Minimal Reflection-Based Environment Management for Go](#minimal-reflection-based-environment-management-for-go)
  - [Getting Started](#getting-started)
  - [Optional \& Defaults](#optional--defaults)
  - [Reading `.env`-Files](#reading-env-files)
  - [Supported Types](#supported-types)
  - [Advanced Usage](#advanced-usage)
      - [Common Prefix](#common-prefix)
      - [Fallback Values](#fallback-values)
      - [Error Handling](#error-handling)

## Getting Started

Using `minienv` is quite simple, just create a struct and annotate it with `env:""` tags:

```go
type Environment struct {
    Value int `env:"VALUE"`
}

var e Environment
if err := minienv.Load(&e); err != nil {
    // handle error
}

print(e.Value) // will equal to whatever the VALUE env variable is set to
```

## Optional & Defaults

The package supports specifying values as optional or providing a default for them. This can be done by using the `optional` and `default=...` keywords in the `env:""`-tag:

```go
type Environment struct {
    Value int `env:"VALUE,optional"` // will be set to the zero value if not provided
    Other string `env:OTHER,default=test"` // will be set to "test" if not provided
}
```

To see the correct syntax for setting default values for map and slice types in the [list of supported types](#supported-types).

## Reading `.env`-Files

`minienv` supports loading variables from `.env` files by using the `WithEnvFile(...)` option:

```go
type Environment struct {
    Port int `env:"PORT"`
}

var e Environment

// the first argument is the path to the env-file and the 
// second controls if it is required or not.
err := minienv.Load(&e, minienv.WithEnvFile(".env", false)) 
if err != nil {
    // handle error
}
```

It is possible to read from multiple `.env`-files by simply providing the option multiple times.

## Supported Types

The following table provides an overview of all supported types and the syntax for default values:

| Type                                            | Default Example            |
| ----------------------------------------------- | -------------------------- |
| `int`, `int8`, `int16`, `int32` & `int64`       | `default=1`                |
| `float32` & `float64`                           | `default=1.1`              |
| `bool`                                          | `default=true`             |
| `string`                                        | `default=example`          |
| `[]...` (Slices support types listed above)     | `default=val\|val\|val`    |
| `map[...]...` (Maps support types listed above) | `default=key:val\|key:val` |


## Advanced Usage

#### Common Prefix

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


#### Fallback Values

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

#### Error Handling

If Minienv encounters any issues during loading, it will raise an error to the enduser. These errors are wrapped in custom error objects that allow you to react to them more precisely.

If the input to the `Load()`-function itself is invalid (e.g. not a pointer), Minienv will raise the predefined `ErrInvalidInput`-error:

```go
var e Environment
err := minienv.Load(e) // e is not a pointer here, therefore invalid
if err == minienv.ErrInvalidInput {
    // do something...
}
```

Additionally if Minienv fails to load a value into a certain field, for example due to a type mismatch, it will raise an error of type `FieldError`:

```go
var e Environment
err := minienv.Load(&e)
if err != nil {
    fieldErr = err.(minienv.FieldError)
    // handle load error
}
```

The `FieldError` additionally exposes the affected field that failed together with the underlying error.