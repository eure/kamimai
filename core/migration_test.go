package core

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kaneshin/kamimai/core/internal/direction"
)

func TestService(t *testing.T) {
	assert := assert.New(t)

	conf := testMustNewConfig(t)
	svc := NewService(conf).WithVersion(123)

	assert.EqualValues(123, svc.version)

	svc.WithVersion(101)
	assert.EqualValues(101, svc.version)

	svc.direction = direction.Up
	svc.apply()
	migs := ([]*Migration)(svc.data)
	assert.EqualValues(1, migs[0].version)
	assert.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_up.sql"))

	svc.direction = direction.Down
	svc.apply()
	migs = ([]*Migration)(svc.data)
	assert.EqualValues(1, migs[0].version)
	assert.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_down.sql"))
}

func TestMigration(t *testing.T) {
	assert := assert.New(t)

	mig := NewMigration().WithVersion(123)
	assert.EqualValues(123, mig.version)

	mig.WithVersion(101)
	assert.EqualValues(101, mig.version)
}

func TestSortMigrations(t *testing.T) {
	assert := assert.New(t)

	migs := (Migrations)([]*Migration{
		NewMigration().WithVersion(123),
		NewMigration().WithVersion(12),
		NewMigration().WithVersion(1023),
		NewMigration().WithVersion(383),
		NewMigration().WithVersion(971),
		NewMigration().WithVersion(184),
	})

	sort.Sort(migs)
	assert.EqualValues(12, migs[0].version)
	assert.EqualValues(123, migs[1].version)
	assert.EqualValues(184, migs[2].version)
	assert.EqualValues(383, migs[3].version)
	assert.EqualValues(971, migs[4].version)
	assert.EqualValues(1023, migs[5].version)
}

func TestMigrationsIndex(t *testing.T) {
	assert := assert.New(t)

	migs := (Migrations)([]*Migration{
		NewMigration().WithVersion(123),
		NewMigration().WithVersion(12),
		NewMigration().WithVersion(1023),
		NewMigration().WithVersion(383),
		NewMigration().WithVersion(971),
		NewMigration().WithVersion(184),
	})
	assert.EqualValues(notFoundIndex, migs.index(Migration{version: 1000000}))

	assert.Equal(1, migs.index(Migration{version: 12}))
	assert.Equal(2, migs.index(Migration{version: 1023}))

	sort.Sort(migs)

	assert.Equal(0, migs.index(Migration{version: 12}))
	assert.Equal(5, migs.index(Migration{version: 1023}))
}
