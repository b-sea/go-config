// Package config is a simple configuration loading helper.
package config

import (
	"fmt"
	"path"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

const (
	defaultDelim    = "."
	defaultEnvDelim = "__"
	defaultTag      = "config"
)

type config struct {
	delim         string
	tag           string
	filePaths     map[string]bool
	envPrefix     string
	envDelim      string
	unmarshalPath string
}

// Load reads and loads configuration data into the given config.
// If multiple files are specified, they are loaded in the order provided.
// Any environment variable overrides will be applied last and will
// always take precedent.
//
// Any configuration structs passed in should use the `config:"..."` struct tags
// for proper unmarshalling by default.
func Load(cfg any, options ...Option) error {
	setup := &config{
		delim:         defaultDelim,
		tag:           defaultTag,
		filePaths:     make(map[string]bool),
		envPrefix:     "",
		envDelim:      defaultEnvDelim,
		unmarshalPath: "",
	}

	for _, option := range options {
		option(setup)
	}

	manager := koanf.New(setup.delim)

	for path := range setup.filePaths {
		if err := loadFile(manager, path); err != nil {
			return err
		}
	}

	err := manager.Load(env.Provider(setup.delim, env.Opt{
		Prefix: setup.envPrefix,
		TransformFunc: func(key string, value string) (string, any) {
			key = strings.ReplaceAll(
				strings.ToLower(strings.TrimPrefix(key, setup.envPrefix)),
				setup.envDelim,
				setup.delim,
			)

			return key, value
		},
	}), nil)
	if err != nil {
		return envError(err) // coverage-ignore
	}

	if err := manager.UnmarshalWithConf(setup.unmarshalPath, cfg, koanf.UnmarshalConf{Tag: setup.tag}); err != nil {
		return unmarshalError(err)
	}

	return nil
}

func loadFile(manager *koanf.Koanf, filepath string) error {
	var parser koanf.Parser

	switch path.Ext(filepath) {
	case ".json":
		parser = json.Parser()
	case ".yaml", ".yml":
		parser = yaml.Parser()
	default:
		return fileLoadError(fmt.Errorf("unknown file type: %s", filepath)) //nolint: err113
	}

	if err := manager.Load(file.Provider(filepath), parser); err != nil {
		return fileLoadError(err)
	}

	return nil
}
