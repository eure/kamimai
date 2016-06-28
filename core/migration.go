package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/kaneshin/kamimai/core/internal/cast"
	"github.com/kaneshin/kamimai/core/internal/direction"
	"github.com/kaneshin/kamimai/core/internal/version"
)

const (
	notFoundIndex = 0xffff
	migrationsDir = "migrations"
)

var (
	errOutOfBoundsMigrations = errors.New("out of bounds migration")
)

type (
	// A Service manages for kamimai.
	Service struct {
		config    *Config
		driver    Driver
		version   uint64
		direction int
		data      Migrations
	}

	// A Migration manages migration files.
	Migration struct {
		version uint64
		name    string
	}

	// A Migrations collects Migration for sorting.
	Migrations []*Migration
)

func (s Service) walker(indexPath map[uint64]*Migration) func(path string, info os.FileInfo, err error) error {
	wd, _ := os.Getwd()

	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		ver := cast.Uint64(version.Get(name))
		mig, found := indexPath[ver]
		if !found || mig == nil {
			mig = &Migration{
				version: ver,
			}
			indexPath[ver] = mig
		}

		if s.direction == direction.Get(name) {
			mig.name = filepath.Clean(filepath.Join(wd, path))
		} else if s.direction == direction.Unknown {
			mig.name = filepath.Clean(filepath.Join(wd, path))
		}

		return nil
	}
}

// NewService returns a new Service pointer that can be chained with builder methods to
// set multiple configuration values inline without using pointers.
func NewService(c *Config) *Service {
	svc := &Service{
		config: c,
	}
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

func (s Service) migrationsDir() string {
	return path.Clean(path.Join(s.config.Dir(), migrationsDir))
}

// MakeMigrationsDir creates a directory named path, along with any necessary
// parents, and returns nil, or else returns an error. If path is
// already a directory, MkdirAll does nothing and returns nil.
func (s *Service) MakeMigrationsDir() error {
	return os.MkdirAll(s.migrationsDir(), 0777)
}

func (s *Service) apply() {
	index := map[uint64]*Migration{}
	filepath.Walk(s.migrationsDir(), s.walker(index))

	list := make([]*Migration, len(index))
	i := 0
	for _, m := range index {
		list[i] = m
		i++
	}

	migs := Migrations(list)
	sort.Sort(migs)
	s.data = migs
}

func (s *Service) do(idx int) error {
	migs := s.data
	if !(0 <= idx && idx < migs.Len()) {
		return errOutOfBoundsMigrations
	}
	return s.driver.Migrate(migs[idx])
}

func (s *Service) step(n int) error {
	s.apply()

	// gets current index of migrations
	idx := s.data.index(Migration{version: s.version})
	if s.version == 0 {
		idx = 0
	}

	// direction of the migration.
	sign := -1
	if n > 0 {
		sign = 1
	}

	for i := 0; i < sign*n; i++ {
		if err := s.do(idx + sign*i); err != nil {
			return err
		}
	}

	return nil
}

// Up upgrades migration version.
func (s *Service) Up() error {
	s.direction = direction.Up
	err := s.step(notFoundIndex)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

// Down downgrades migration version.
func (s *Service) Down() error {
	s.direction = direction.Down
	err := s.step(-notFoundIndex)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

// Next upgrades migration version.
func (s *Service) Next() error {
	s.direction = direction.Up
	return s.step(1)
}

// Prev downgrades migration version.
func (s *Service) Prev() error {
	s.direction = direction.Down
	return s.step(-1)
}

// NextMigration returns next version migrations.
func (s *Service) NextMigration(name string) (up *Migration, down *Migration, err error) {
	s.apply()

	// initialize default variables for making migrations.
	up, down = &Migration{version: 1, name: ""}, &Migration{version: 1, name: ""}
	ver := "001"

	// gets the oldest migration version file.
	if obj := s.data.first(); obj != nil {
		// for version format
		_, file := filepath.Split(obj.name)
		ver = version.Format(file)
	}
	if obj := s.data.last(); obj != nil {
		// for version number
		v := obj.version + 1
		up.version, down.version = v, v
	}

	// [ver]_[name]_[direction-suffix][.ext]
	base := fmt.Sprintf("%s_%s_%%s%%s", ver, name)
	// including dot
	ext := s.driver.Ext()

	// up
	n := fmt.Sprintf(base, up.version, direction.Suffix(direction.Up), ext)
	up.name = filepath.Join(s.migrationsDir(), n)
	// down
	n = fmt.Sprintf(base, down.version, direction.Suffix(direction.Down), ext)
	down.name = filepath.Join(s.migrationsDir(), n)

	return
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

// Version returns a migration version value.
func (m Migration) Version() uint64 {
	return m.version
}

// Read reads from file until an error or EOF and returns the data it read.
// A successful call returns err == nil, not err == EOF. Because ReadAll is
// defined to read from src until EOF, it does not treat an EOF from Read
// as an error to be reported.
func (m Migration) Read() ([]byte, error) {
	file, err := os.Open(m.name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

// Name returns a file name.
func (m Migration) Name() string {
	return m.name
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

func (m Migrations) first() *Migration {
	if m.Len() == 0 {
		return nil
	}
	return m[0]
}

func (m Migrations) last() *Migration {
	c := m.Len()
	if c == 0 {
		return nil
	}
	return m[c-1]
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
