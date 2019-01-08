package kamimai

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mr04vv/kamimai/core"
	_ "github.com/mr04vv/kamimai/driver"
	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {
	assert := assert.New(t)

	conf, err := core.NewConfig("examples/testdata")
	if assert.NoError(err) {
		conf.WithEnv("development")

		defer func() {
			// Down
			assert.NoError(Down(conf))
			ver, err := Current(conf)
			if assert.NoError(err) {
				assert.EqualValues(0, ver)
			}
		}()

		// Down
		assert.NoError(Down(conf))
		ver, err := Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(0, ver)
		}

		// Sync
		assert.NoError(Sync(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(1, ver)
		}

		// Down
		assert.NoError(Down(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(0, ver)
		}

		// Sync
		assert.NoError(Sync(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(1, ver)
		}
	}
}

func TestUp(t *testing.T) {
	assert := assert.New(t)

	conf, err := core.NewConfig("examples/testdata")
	if assert.NoError(err) {
		conf.WithEnv("development")

		defer func() {
			// Down
			assert.NoError(Down(conf))
			ver, err := Current(conf)
			if assert.NoError(err) {
				assert.EqualValues(0, ver)
			}
		}()

		// Down
		assert.NoError(Down(conf))
		ver, err := Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(0, ver)
		}

		// Up
		assert.NoError(Up(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(1, ver)
		}

		// Down
		assert.NoError(Down(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(0, ver)
		}

		// Up
		assert.NoError(Up(conf))
		ver, err = Current(conf)
		if assert.NoError(err) {
			assert.EqualValues(1, ver)
		}
	}
}
