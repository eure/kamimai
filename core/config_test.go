package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	os.Clearenv()
	os.Setenv("MYSQL_USER", "testuser")
	os.Setenv("MYSQL_PASSWORD", "testpassword")
	os.Setenv("MYSQL_DATABASE", "kamimai")
}

func testMustNewConfig(t *testing.T) *Config {
	conf, err := NewConfig("../examples/testdata")
	if assert.NoError(t, err) {
		if assert.NotNil(t, conf) {
			conf.WithEnv("development")
		}
	}
	return conf
}

func TestNewConfig(t *testing.T) {
	asrt := assert.New(t)

	var (
		conf *Config
		err  error
	)

	conf, err = NewConfig("")
	asrt.Nil(conf)
	asrt.Error(err)

	asrt.Equal("", Config{}.Driver())
	asrt.Equal("", Config{}.Dsn())
	asrt.Equal("", Config{}.migrationsDir())

	conf, err = NewConfig("../examples/testdata")
	asrt.NotNil(conf)

	conf.WithEnv("development")
	if asrt.NoError(err) {
		asrt.Equal("mysql", conf.Driver())
		asrt.Equal("mysql://testuser:testpassword@tcp(:)/kamimai?charset=utf8", conf.Dsn())
		asrt.Equal("../examples/testdata/migrations", conf.migrationsDir())
	}

	conf.WithEnv("test")
	if asrt.NoError(err) {
		asrt.Equal("mysql", conf.Driver())
		asrt.Equal("mysql://testuser:testpassword@tcp(:)/kamimai?charset=utf8", conf.Dsn())
		asrt.Equal("../examples/testdata/test", conf.migrationsDir())
	}

	var (
		confMySQL  *Config
		confSQLite *Config
	)

	confMySQL, err = NewConfig("../examples/mysql")
	confMySQL.WithEnv("development")
	if asrt.NoError(err) {
		asrt.Equal("mysql", confMySQL.Driver())
	}

	confSQLite, err = NewConfig("../examples/sqlite3")
	confSQLite.WithEnv("development")
	if asrt.NoError(err) {
		asrt.Equal("sqlite3", confSQLite.Driver())
	}

	conf = MergeConfig(confMySQL, confSQLite)
	conf.WithEnv("development")
	if asrt.NoError(err) {
		asrt.Equal("mysql", conf.Driver())
	}

	conf = MergeConfig(confSQLite, confMySQL)
	conf.WithEnv("development")
	if asrt.NoError(err) {
		asrt.Equal("sqlite3", conf.Driver())
	}
}
