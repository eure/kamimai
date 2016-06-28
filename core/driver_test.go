package core

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type drv struct {
}

func (d *drv) Open(_ string) error                     { return nil }
func (d *drv) Close() error                            { return nil }
func (d *drv) Ext() string                             { return "" }
func (d *drv) Transaction(f func(*sql.Tx) error) error { return nil }
func (d *drv) Migrate(mig *Migration) error            { return nil }
func (d *drv) Version() Version                        { return nil }

func TestRegisterDriver(t *testing.T) {
	assert := assert.New(t)

	assert.NotNil(registry)

	assert.Panics(func() {
		RegisterDriver("nil", (Driver)(nil))
	})

	// gets and sets
	assert.Nil(GetDriver("driver"), "should be nil.")
	RegisterDriver("driver", &drv{})
	assert.NotNil(GetDriver("driver"), "should be retrieve a registered driver.")

	assert.Panics(func() {
		RegisterDriver("driver", &drv{})
	}, "should be panic if registering driver twice.")
}
