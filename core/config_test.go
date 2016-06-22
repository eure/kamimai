package core

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)
	os.Clearenv()
	os.Setenv("MYSQL_USER", "testuser")
	os.Setenv("MYSQL_PASSWORD", "testpassword")

	var (
		conf *Config
		err  error
	)

	conf, err = NewConfig("")
	assert.Nil(conf)
	assert.Error(err)

	assert.Equal("", Config{}.Import())
	assert.Equal("", Config{}.Driver())
	assert.Equal("", Config{}.Dsn())

	conf, err = NewConfig("../examples/mysql")
	assert.NotNil(conf)
	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", conf.Driver())
		assert.Equal("github.com/go-sql-driver/mysql", conf.Import())
		assert.Equal("mysql://testuser:testpassword@127.0.0.1:3306/kamimai?charset=utf8", conf.Dsn())
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
