package config_test

import (
	"os"
	"testing"

	"github.com/b-sea/go-config/config"
	"github.com/stretchr/testify/assert"
)

type Child struct {
	Data any `config:"data"`
}

type Config struct {
	Data  any   `config:"data"`
	Child Child `config:"child"`
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	type testCase struct {
		options []config.Option
		env     map[string]string
		input   any
		result  any
		err     error
	}

	tests := map[string]testCase{
		"load file": {
			options: []config.Option{
				config.WithFile("./../samples/simple.yml"),
			},
			env: map[string]string{},
			input: Config{
				Data: "initial value",
				Child: Child{
					Data: "initial value",
				},
			},
			result: Config{
				Data: "simple value",
				Child: Child{
					Data: "simple value",
				},
			},
			err: nil,
		},
		"load multiple files": {
			options: []config.Option{
				config.WithFile("./../samples/partial.yml"),
				config.WithFile("./../samples/partial.json"),
			},
			env: map[string]string{},
			input: Config{
				Data: "initial value",
				Child: Child{
					Data: "initial value",
				},
			},
			result: Config{
				Data: "json value",
				Child: Child{
					Data: "yml value",
				},
			},
			err: nil,
		},
		"unsupported file type": {
			options: []config.Option{
				config.WithFile("./../samples/simple.ini"),
			},
			env:    map[string]string{},
			input:  Config{},
			result: Config{},
			err:    config.ErrFileLoad,
		},
		"missing file": {
			options: []config.Option{
				config.WithFile("./../samples/random.yml"),
			},
			env:    map[string]string{},
			input:  Config{},
			result: Config{},
			err:    config.ErrFileLoad,
		},
		"env override": {
			options: []config.Option{
				config.WithFile("./../samples/simple.yml"),
				config.WithEnvPrefix("CONFIG_UNITTEST"),
			},
			env: map[string]string{
				"CONFIG_UNITTEST_DATA": "env value",
			},
			input: Config{
				Data: "initial value",
				Child: Child{
					Data: "initial value",
				},
			},
			result: Config{
				Data: "env value",
				Child: Child{
					Data: "simple value",
				},
			},
			err: nil,
		},
		"custom env delim": {
			options: []config.Option{
				config.WithFile("./../samples/simple.yml"),
				config.WithEnvPrefix("CONFIG_UNITTEST"),
				config.SetEnvDelim("_S_"),
			},
			env: map[string]string{
				"CONFIG_UNITTEST_CHILD_S_DATA": "env value",
			},
			input: Config{
				Data: "initial value",
				Child: Child{
					Data: "initial value",
				},
			},
			result: Config{
				Data: "simple value",
				Child: Child{
					Data: "env value",
				},
			},
			err: nil,
		},
		"unmarshal path": {
			options: []config.Option{
				config.WithFile("./../samples/simple.yml"),
				config.WithUnmarshalPath("child"),
			},
			env: map[string]string{},
			input: Child{
				Data: "initial value",
			},
			result: Child{
				Data: "simple value",
			},
			err: nil,
		},
		"unmarshal error": {
			options: []config.Option{
				config.WithFile("./../samples/simple.yml"),
			},
			env:    map[string]string{},
			input:  "something",
			result: 123,
			err:    config.ErrUnmarshal,
		},
		"custom tag": {
			options: []config.Option{
				config.WithFile("./../samples/partial.json"),
				config.SetTag("special"),
			},
			env: map[string]string{},
			input: struct {
				Data string `special:"data"`
			}{
				Data: "initial",
			},
			result: struct {
				Data string `special:"data"`
			}{
				Data: "json value",
			},
			err: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			for k, v := range test.env {
				os.Setenv(k, v)
			}

			err := config.Load(&test.input, test.options...)

			for k := range test.env {
				os.Unsetenv(k)
			}

			if test.err == nil {
				assert.Equal(t, test.result, test.input)
				assert.NoError(t, err)
			} else {
				assert.ErrorAs(t, err, &test.err)
			}
		})
	}
}
