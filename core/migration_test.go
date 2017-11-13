package core

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration(t *testing.T) {
	asrt := assert.New(t)

	mig := NewMigration().WithVersion(123)
	asrt.EqualValues(123, mig.version)

	mig.WithVersion(101)
	asrt.EqualValues(101, mig.version)
}

func TestSortMigrations(t *testing.T) {
	asrt := assert.New(t)

	migs := (Migrations)([]*Migration{
		NewMigration().WithVersion(123),
		NewMigration().WithVersion(12),
		NewMigration().WithVersion(1023),
		NewMigration().WithVersion(383),
		NewMigration().WithVersion(971),
		NewMigration().WithVersion(184),
	})

	sort.Sort(migs)
	asrt.EqualValues(12, migs[0].version)
	asrt.EqualValues(123, migs[1].version)
	asrt.EqualValues(184, migs[2].version)
	asrt.EqualValues(383, migs[3].version)
	asrt.EqualValues(971, migs[4].version)
	asrt.EqualValues(1023, migs[5].version)
}

func TestMigrationsIndex(t *testing.T) {
	asrt := assert.New(t)

	migs := (Migrations)([]*Migration{
		NewMigration().WithVersion(123),
		NewMigration().WithVersion(12),
		NewMigration().WithVersion(1023),
		NewMigration().WithVersion(383),
		NewMigration().WithVersion(971),
		NewMigration().WithVersion(184),
	})
	asrt.EqualValues(notFoundIndex, migs.index(Migration{version: 1000000}))

	asrt.Equal(1, migs.index(Migration{version: 12}))
	asrt.Equal(2, migs.index(Migration{version: 1023}))

	sort.Sort(migs)

	asrt.Equal(0, migs.index(Migration{version: 12}))
	asrt.Equal(5, migs.index(Migration{version: 1023}))
}
