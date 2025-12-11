# go-config

A small package to make loading configration data simpler


## Install

```bash
go get github.com/b-sea/go-config
```

## Load from File

One or more config files can be loaded into the given config struct.

```go
package main

import (
    "os"

    "github.com/b-sea/go-config/config"
)

func main() {
    var cfg struct {
        Foo int `config:"foo"`
        Bar string `config:"bar"`
        Nested struct {
            Baz string `config:"baz`
        } `config:"nested"`
    }
    
    if err := config.Load(&cfg, config.WithFile("/my/config/file.yml")); err != nil {
        os.Exit(1)
    }

    ...
}
```

### Multiple Config Files

When multiple files are used, they are loaded in the given order and are able to build off of each other. Any shared 
keys will be overwritten as they are discovered.

```go
package main

import (
    "os"

    "github.com/b-sea/go-config/config"
)

func main() {
    var cfg struct {
        Foo int `config:"foo"`
        Bar string `config:"bar"`
        Nested struct {
            Baz string `config:"baz`
        } `config:"nested"`
    }
    
    err := config.Load(
        &cfg, 
        config.WithFile("/my/config/primary.yml"), 
        config.WithFile("/my/config/secondary.json"), 
        config.WithFile("/my/config/teriary.yml"),
    )
    if err != nil {
        os.Exit(1)
    }

    ...
}
```

## Environment Variable Overrides

Environment variables are loaded after any files. If no files are provided, only environment variables will be searched. To denote value hierarchy, `__` should be used as the delimiter. This can be controlled with the `SetEnvDelim` option.

```go
package main

import (
    "os"

    "github.com/b-sea/go-config/config"
)

func main() {
    var cfg struct {
        Foo int `config:"foo"`
        Bar string `config:"bar"`
        Nested struct {
            Baz string `config:"baz`
        } `config:"nested"`
    }

    os.Setenv("FOO", "123")
    os.Setenv("NESTED__BAZ", "some value")
    
    err := config.Load(&cfg)
    if err != nil {
        os.Exit(1)
    }

    ...
}
```

### Environment Prefix

An environment variable prefix can be set with the `WithEnvPrefix` option

```go
package main

import (
    "os"

    "github.com/b-sea/go-config/config"
)

func main() {
    var cfg struct {
        Foo int `config:"foo"`
        Bar string `config:"bar"`
        Nested struct {
            Baz string `config:"baz`
        } `config:"nested"`
    }

    os.Setenv("FOO", "123") // No prefix, will not be loaded
    os.Setenv("MY_PREFIX_NESTED__BAZ", "some value") // Will be loaded
    
    err := config.Load(&cfg, config.WithEnvPrefix("MY_PREFIX"))
    if err != nil {
        os.Exit(1)
    }

    ...
}
```

## Custom Tag

By default, the `config` struct tag will be used to match fields to parsed values. The `SetTag` option can be used to override this.

```go
package main

import (
    "os"

    "github.com/b-sea/go-config/config"
)

func main() {
    var cfg struct {
        Foo int `special:"foo"`
        Bar string `special:"bar"`
        Nested struct {
            Baz string `special:"baz`
        } `special:"nested"`
    }
 
    err := config.Load(&cfg, config.SetTag("special"))
    if err != nil {
        os.Exit(1)
    }

    ...
}
```
