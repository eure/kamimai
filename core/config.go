package core

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type (
	// Config object.
	Config struct {
		Data map[string]Internal
		env  string
		dir  string
	}

	Internal struct {
		Driver    string `yaml:"driver"`
		Dsn       string `yaml:"dsn"`
		Directory string `yaml:"directory"`
	}
)

// MustNewConfig returns a new config. dir cannot be empty.
func MustNewConfig(dir string) *Config {
	c, err := NewConfig(dir)
	if err != nil {
		panic(err)
	}
	return c
}

// NewConfig returns a new config. dir cannot be empty.
func NewConfig(dir string) (*Config, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	conf := Config{dir: dir}
	for _, f := range files {
		fpath := filepath.Join(dir, f.Name())

		switch filepath.Ext(fpath) {
		case ".yml", ".yaml":
			b, err := ioutil.ReadFile(fpath)
			if err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(b, &conf.Data); err != nil {
				return nil, err
			}

		case ".tml", ".toml":
			if _, err := toml.DecodeFile(fpath, &conf.Data); err != nil {
				return nil, err
			}
		}
	}

	return &conf, nil
}

// MergeConfig returns merged config.
func MergeConfig(c ...*Config) *Config {
	conf := Config{Data: map[string]Internal{}}
	for _, v := range c {
		if conf.env == "" && v.env != "" {
			conf.env = v.env
		}

		for key, vv := range v.Data {
			if _, ok := conf.Data[key]; !ok {
				conf.Data[key] = vv
			}
		}
	}
	return &conf
}

// WithEnv sets an environment of config.
func (c *Config) WithEnv(env string) *Config {
	c.env = env
	return c
}

// Dir returns a config file existing path name.
func (c Config) Dir() string {
	return c.dir
}

// Driver returns a raw driver string.
func (c Config) Driver() string {
	if d, ok := c.Data[c.env]; ok {
		return d.Driver
	}
	return ""
}

// Dsn returns a raw dsn string.
func (c Config) Dsn() string {
	if d, ok := c.Data[c.env]; ok {
		return os.ExpandEnv(d.Dsn)
	}
	return ""
}

func (c Config) migrationsDir() string {
	if d, ok := c.Data[c.env]; ok {
		if d.Directory == "" {
			d.Directory = "migrations"
		}
		return path.Clean(path.Join(c.Dir(), d.Directory))
	}
	return ""
}
