package core

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	assert := assert.New(t)

	mig := NewMigration(nil).WithVersion(123)
	assert.EqualValues(123, mig.version)

	mig.WithVersion(101)
	assert.EqualValues(101, mig.version)
}

func TestSortMigrations(t *testing.T) {
	assert := assert.New(t)

	migs := (Migrations)([]*Migration{
		NewMigration(nil).WithVersion(123),
		NewMigration(nil).WithVersion(12),
		NewMigration(nil).WithVersion(1023),
		NewMigration(nil).WithVersion(383),
		NewMigration(nil).WithVersion(971),
		NewMigration(nil).WithVersion(184),
	})

	sort.Sort(migs)
	assert.EqualValues(12, migs[0].version)
	assert.EqualValues(123, migs[1].version)
	assert.EqualValues(184, migs[2].version)
	assert.EqualValues(383, migs[3].version)
	assert.EqualValues(971, migs[4].version)
	assert.EqualValues(1023, migs[5].version)
}
