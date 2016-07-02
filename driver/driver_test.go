package driver

import (
	"testing"

	"github.com/kaneshin/kamimai/core"
	"github.com/stretchr/testify/assert"
)

type harness struct {
	driver core.Driver
	dsn    string
}

func testDriver(t *testing.T, h harness) {
	assert := assert.New(t)

	if assert.NoError(h.driver.Open(h.dsn)) {
		testVersion(t, h.driver.Version())

		assert.NoError(h.driver.Close())
	}
}

func testVersion(t *testing.T, version core.Version) {
	assert := assert.New(t)
	var (
		err error
		val uint64
	)

	if assert.NoError(version.Drop()) {
		// current
		val, err = version.Current()
		if assert.Error(err) {
			assert.EqualValues(0, val)
		}
		assert.NoError(version.Create())
	}

	// current
	val, err = version.Current()
	if assert.NoError(err) {
		assert.EqualValues(0, val)
	}

	assert.NoError(version.Insert(1))
	assert.EqualValues(0, version.Count(100))
	assert.NoError(version.Insert(100))
	val, err = version.Current()
	if assert.NoError(err) {
		assert.EqualValues(100, val, "should be 100")
		assert.EqualValues(1, version.Count(100))

		// delete
		assert.NoError(version.Delete(50))
		val, err = version.Current()
		if assert.NoError(err) {
			assert.EqualValues(100, val, "should be 100")

			// delete
			assert.NoError(version.Delete(100))
			val, err = version.Current()
			if assert.NoError(err) {
				assert.EqualValues(1, val, "should be 1")
				assert.EqualValues(0, version.Count(100))
			}
		}
	}
}

func TestMySQLDriver(t *testing.T) {
	assert := assert.New(t)

	driver := new(MySQL)
	assert.Implements((*core.Driver)(nil), driver)
	assert.Implements((*core.Version)(nil), driver.Version())

	conf, err := core.NewConfig("../examples/testdata")
	if assert.NoError(err) {
		conf.WithEnv("development")
		testDriver(t, harness{driver, conf.Dsn()})
	}
}
