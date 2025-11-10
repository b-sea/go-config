package config

import (
	"errors"
	"fmt"
)

// ErrFileLoad, et al. are the custom errors thrown while parsing and loading configs.
var (
	ErrFileLoad  = errors.New("config file load error")
	ErrEnvLoad   = errors.New("config env load error")
	ErrUnmarshal = errors.New("config unmarshal error")
)

func fileLoadError(err error) error {
	return fmt.Errorf("%w: %w", ErrFileLoad, err)
}

func envError(err error) error {
	return fmt.Errorf("%w: %w", ErrEnvLoad, err)
}

func unmarshalError(err error) error {
	return fmt.Errorf("%w: %w", ErrUnmarshal, err)
}
