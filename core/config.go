package core

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v2"
)

type (
	// Config object.
	Config struct {
		data map[string]internal
		env  string
	}

	internal struct {
		Driver string `yaml:"driver"`
		Dsn    string `yaml:"dsn"`
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

	conf := Config{}
	for _, f := range files {
		fpath := filepath.Join(dir, f.Name())

		switch filepath.Ext(fpath) {
		case ".yml", ".yaml":
			b, err := ioutil.ReadFile(fpath)
			if err != nil {
				return nil, err
			}
			if err := yaml.Unmarshal(b, &conf.data); err != nil {
				return nil, err
			}

		case ".tml", ".toml":
			if _, err := toml.DecodeFile(fpath, &conf.data); err != nil {
				return nil, err
			}
		}
	}

	return &conf, nil
}

// MergeConfig returns merged config.
func MergeConfig(c ...*Config) *Config {
	conf := Config{data: map[string]internal{}}
	for _, v := range c {
		if conf.env == "" && v.env != "" {
			conf.env = v.env
		}

		for key, vv := range v.data {
			if _, ok := conf.data[key]; !ok {
				conf.data[key] = vv
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

// Import returns an import path.
func (c Config) Import() string {
	switch c.Driver() {
	case "mymysql":
		return "github.com/ziutek/mymysql/godrv"
	case "mysql":
		return "github.com/go-sql-driver/mysql"
	case "sqlite3", "sqlite":
		return "github.com/mattn/go-sqlite3"
	}
	return ""
}

// Driver returns a raw driver string.
func (c Config) Driver() string {
	if d, ok := c.data[c.env]; ok {
		return d.Driver
	}
	return ""
}

// Dsn returns a raw dsn string.
func (c Config) Dsn() string {
	if d, ok := c.data[c.env]; ok {
		return os.ExpandEnv(d.Dsn)
	}
	return ""
}
