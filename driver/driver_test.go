package driver

import (
	"database/sql"
	"errors"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mr04vv/kamimai/core"
	"github.com/stretchr/testify/assert"
)

type harness struct {
	driver core.Driver
	dsn    string
}

func testDriver(t *testing.T, h harness) {
	asrt := assert.New(t)

	if asrt.NoError(h.driver.Open(h.dsn)) {
		testVersion(t, h.driver.Version())

		asrt.NoError(h.driver.Transaction(func(tx *sql.Tx) error {
			return nil
		}))

		asrt.Error(h.driver.Transaction(func(tx *sql.Tx) error {
			return errors.New("Error")
		}))

		asrt.NoError(h.driver.Close())
	}
}

func testVersion(t *testing.T, version core.Version) {
	asrt := assert.New(t)
	var (
		err error
		val uint64
	)

	if asrt.NoError(version.Drop()) {
		// current
		val, err = version.Current()
		if asrt.Error(err) {
			asrt.EqualValues(0, val)
		}
		asrt.NoError(version.Create())
	}

	// current
	val, err = version.Current()
	if asrt.NoError(err) {
		asrt.EqualValues(0, val)
	}

	asrt.NoError(version.Insert(1))
	asrt.EqualValues(0, version.Count(100))
	asrt.NoError(version.Insert(100))
	val, err = version.Current()
	if asrt.NoError(err) {
		asrt.EqualValues(100, val, "should be 100")
		asrt.EqualValues(1, version.Count(100))

		// delete
		asrt.NoError(version.Delete(50))
		val, err = version.Current()
		if asrt.NoError(err) {
			asrt.EqualValues(100, val, "should be 100")

			// delete
			asrt.NoError(version.Delete(100))
			val, err = version.Current()
			if asrt.NoError(err) {
				asrt.EqualValues(1, val, "should be 1")
				asrt.EqualValues(0, version.Count(100))
			}
		}
	}
}

func TestMySQLDriver(t *testing.T) {
	asrt := assert.New(t)

	driver := new(MySQL)
	asrt.Implements((*core.Driver)(nil), driver)
	asrt.Implements((*core.Version)(nil), driver.Version())

	conf, err := core.NewConfig("../examples/testdata")
	if asrt.NoError(err) {
		conf.WithEnv("development")
		testDriver(t, harness{driver, conf.Dsn()})
	}

	asrt.Equal(".sql", driver.Ext())
}
