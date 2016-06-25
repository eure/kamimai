package core

import (
	"github.com/kaneshin/kamimai/core/internal/cast"
	"sort"
)

const (
	notFoundIndex = 0xffff
)

type (
	// A Service manages for kamimai.
	Service struct {
		version uint64
		config  *Config
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
	svc := &Service{config: c}

	var migs Migrations
	sort.Sort(migs)
	svc.data = migs

	return svc
}

// WithVersion sets a config version value returning a Config pointer
// for chaining.
func (s *Service) WithVersion(v interface{}) *Service {
	s.version = cast.Uint64(v)
	return s
}

func (s Service) step(n int) error {
	idx := s.data.index(Migration{version: s.version})
	_ = idx
	return nil
}

// Up upgrades migration version.
func (s Service) Up() error {
	return nil
}

// Down downgrades migration version.
func (s Service) Down() error {
	return nil
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
