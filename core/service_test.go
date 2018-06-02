package core

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eure/kamimai/internal/direction"
)

func TestService(t *testing.T) {
	assrt := assert.New(t)

	conf := testMustNewConfig(t)
	svc := NewService(conf).WithVersion(123)

	assrt.EqualValues(123, svc.version)

	svc.WithVersion(101)
	assrt.EqualValues(101, svc.version)

	svc.direction = direction.Up
	assrt.NoError(svc.apply())
	migs := ([]*Migration)(svc.data)
	assrt.EqualValues(1, migs[0].version)
	assrt.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_up.sql"))

	svc.direction = direction.Down
	assrt.NoError(svc.apply())
	migs = ([]*Migration)(svc.data)
	assrt.EqualValues(1, migs[0].version)
	assrt.True(strings.HasSuffix(migs[0].name, "migrations/001_create_product_down.sql"))
}

func TestInvalidService(t *testing.T) {
	assrt := assert.New(t)

	conf := MustNewConfig("../examples/invalid")
	conf.WithEnv("development")

	svc := NewService(conf)

	svc.direction = direction.Up
	err := svc.apply()
	assrt.Error(err)
}
