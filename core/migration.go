package core

import (
	"errors"
	"sort"

	"github.com/kaneshin/kamimai/core/internal/cast"
)

const (
	notFoundIndex = 0xffff
)

var (
	errOutOfBoundsMigrations = errors.New("out of bounds migration")
)

type (
	// A Service manages for kamimai.
	Service struct {
		config  *Config
		driver  Driver
		version uint64
		data    Migrations
	}

	// A Migration manages migration files.
	Migration struct {
		version uint64
	}

	// A Migrations collects Migration for sorting.
	Migrations []*Migration
)

// NewService returns a new Service pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
func NewService(c *Config) *Service {
	svc := &Service{
		config: c,
	}
	// TODO:
	var migs Migrations
	sort.Sort(migs)
	svc.data = migs
	// //-TODO:
	return svc
}

// WithVersion sets a config version value returning a Service pointer
// for chaining.
func (s *Service) WithVersion(v interface{}) *Service {
	s.version = cast.Uint64(v)
	return s
}

// WithDriver sets a driver returning a Service pointer for chaining.
func (s *Service) WithDriver(d Driver) *Service {
	s.driver = d
	return s
}

func (s Service) do(idx int) error {
	migs := s.data
	if !(0 <= idx && idx < migs.Len()) {
		return errOutOfBoundsMigrations
	}
	return s.driver.Migrate(migs[idx])
}

func (s Service) step(n int) error {
	// gets current index of migrations
	idx := s.data.index(Migration{version: s.version})

	// direction of migrations.
	if n > 0 {
		for i := 1; i <= n; i++ {
			if err := s.do(idx + i); err != nil {
				return err
			}
		}
	} else {
		for i := 1; i <= -n; i++ {
			if err := s.do(idx - i); err != nil {
				return err
			}
		}
	}

	return nil
}

// Up upgrades migration version.
func (s Service) Up() error {
	err := s.step(notFoundIndex)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

// Down downgrades migration version.
func (s Service) Down() error {
	err := s.step(-notFoundIndex)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

// Next upgrades migration version.
func (s Service) Next() error {
	return s.step(1)
}

// Prev downgrades migration version.
func (s Service) Prev() error {
	return s.step(-1)
}

//////////////////////////////
// Migration
//////////////////////////////

// NewMigration returns a new Migration pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
func NewMigration() *Migration {
	return &Migration{}
}

// WithVersion sets a config version value returning a Config pointer
// for chaining.
func (m *Migration) WithVersion(v interface{}) *Migration {
	m.version = cast.Uint64(v)
	return m
}

//////////////////////////////
// Migrations
//////////////////////////////

func (m Migrations) index(mig Migration) int {
	for i, v := range m {
		if v.version == mig.version {
			return i
		}
	}
	return int(notFoundIndex)
}

// Len is the number of elements in the collection.
// Required by Sort Interface{}
func (m Migrations) Len() int {
	return len(m)
}

// Less reports whether the element with
// index i should sort before the element with index j.
// Required by Sort Interface{}
func (m Migrations) Less(i, j int) bool {
	return m[i].version < m[j].version
}

// Swap swaps the elements with indexes i and j.
// Required by Sort Interface{}
func (m Migrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
