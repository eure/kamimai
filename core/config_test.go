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
	assert := assert.New(t)

	var (
		conf *Config
		err  error
	)

	conf, err = NewConfig("")
	assert.Nil(conf)
	assert.Error(err)

	assert.Equal("", Config{}.Driver())
	assert.Equal("", Config{}.Dsn())
	assert.Equal("", Config{}.migrationsDir())

	conf, err = NewConfig("../examples/testdata")
	assert.NotNil(conf)

	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", conf.Driver())
		assert.Equal("mysql://testuser:testpassword@tcp(:)/kamimai?charset=utf8", conf.Dsn())
		assert.Equal("../examples/testdata/migrations", conf.migrationsDir())
	}

	conf.WithEnv("test")
	if assert.NoError(err) {
		assert.Equal("mysql", conf.Driver())
		assert.Equal("mysql://testuser:testpassword@tcp(:)/kamimai?charset=utf8", conf.Dsn())
		assert.Equal("../examples/testdata/test", conf.migrationsDir())
	}

	var (
		confMySQL  *Config
		confSQLite *Config
	)

	confMySQL, err = NewConfig("../examples/mysql")
	confMySQL.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", confMySQL.Driver())
	}

	confSQLite, err = NewConfig("../examples/sqlite3")
	confSQLite.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("sqlite3", confSQLite.Driver())
	}

	conf = MergeConfig(confMySQL, confSQLite)
	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", conf.Driver())
	}

	conf = MergeConfig(confSQLite, confMySQL)
	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("sqlite3", conf.Driver())
	}
}
