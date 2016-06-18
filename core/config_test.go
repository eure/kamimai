package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	assert := assert.New(t)

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
		assert.Equal("mysql://$DB_USER:$DB_PASSWD@127.0.0.1:3306/kamimai?charset=utf8&keepalive=1200", conf.Dsn())
	}

	var (
		conf1 *Config
		conf2 *Config
		conf3 *Config
	)

	conf1, err = NewConfig("../examples/mysql")
	conf1.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", conf1.Driver())
	}

	conf2, err = NewConfig("../examples/mymysql")
	conf2.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mymysql", conf2.Driver())
	}

	conf3, err = NewConfig("../examples/sqlite3")
	conf3.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("sqlite3", conf3.Driver())
	}

	conf = MergeConfig(conf1, conf2, conf3)
	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mysql", conf.Driver())
	}

	conf = MergeConfig(conf2, conf1, conf3)
	conf.WithEnv("development")
	if assert.NoError(err) {
		assert.Equal("mymysql", conf.Driver())
	}
}
