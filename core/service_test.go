package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eure/kamimai/internal/direction"
)

func TestService(t *testing.T) {
	assert := assert.New(t)

	conf := testMustNewConfig(t)
	svc := NewService(conf).WithVersion(123)

	assert.EqualValues(123, svc.version)

	svc.WithVersion(101)
	assert.EqualValues(101, svc.version)

	svc.direction = direction.Up
	assert.NoError(svc.apply())
	migs := ([]*Migration)(svc.data)
	assert.EqualValues(1, migs[0].version)
	assert.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_up.sql"))

	svc.direction = direction.Down
	assert.NoError(svc.apply())
	migs = ([]*Migration)(svc.data)
	assert.EqualValues(1, migs[0].version)
	assert.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_down.sql"))
}

func TestInvalidService(t *testing.T) {

	assert := assert.New(t)

	conf := MustNewConfig("../examples/invalid")
	conf.WithEnv("development")

	svc := NewService(conf)

	svc.direction = direction.Up
	err := svc.apply()
	assert.Error(err)
}
