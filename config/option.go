package config

import "strings"

// Option is a configuration loading option.
type Option func(c *config)

// SetTag sets the struct tag to unmarshal.
func SetTag(tag string) Option {
	return func(c *config) {
		if tag == "" {
			return
		}

		c.tag = tag
	}
}

// WithUnmarshalPath sets the struct key path to unmarshal.
func WithUnmarshalPath(path string) Option {
	return func(c *config) {
		c.unmarshalPath = path
	}
}

// WithEnvPrefix sets the environment variable prefix to use for config overrides.
func WithEnvPrefix(prefix string) Option {
	return func(c *config) {
		if prefix != "" && !strings.HasSuffix(prefix, "_") {
			prefix += "_"
		}

		c.envPrefix = prefix
	}
}

// SetDelim sets the config hierarch delimiter.
func SetDelim(delim string) Option {
	return func(c *config) {
		if delim == "" {
			return
		}

		c.delim = delim
	}
}

// SetEnvDelim sets the environment variable hierarchy delimiter.
func SetEnvDelim(delim string) Option {
	return func(c *config) {
		if delim == "" {
			return
		}

		c.envDelim = strings.ToLower(delim)
	}
}

// WithFile adds a configuration file to load.
func WithFile(path string) Option {
	return func(c *config) {
		c.filePaths[path] = true
	}
}
