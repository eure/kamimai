package core

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/eure/kamimai/internal/cast"
	"github.com/eure/kamimai/internal/direction"
	"github.com/eure/kamimai/internal/version"
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
		config    *Config
		driver    Driver
		version   uint64
		direction int
		data      Migrations
	}
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

// MakeMigrationsDir creates a directory named path, along with any necessary
// parents, and returns nil, or else returns an error. If path is
// already a directory, MkdirAll does nothing and returns nil.
func (s *Service) MakeMigrationsDir() error {
	return os.MkdirAll(s.config.migrationsDir(), 0777)
}

func (s *Service) apply() {
	index := map[uint64]*Migration{}
	filepath.Walk(s.config.migrationsDir(), s.walker(index))

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
	drv := s.driver

	migs := s.data
	if !(0 <= idx && idx < migs.Len()) {
		return errOutOfBoundsMigrations
	}

	mig := migs[idx]
	if err := drv.Migrate(mig); err != nil {
		return err
	}

	switch s.direction {
	case direction.Up:
		if err := drv.Version().Insert(mig.version); err != nil {
			return err
		}
	case direction.Down:
		if err := drv.Version().Delete(mig.version); err != nil {
			return err
		}
	}

	log.Println("applied", mig.Name())

	return nil
}

func (s *Service) step(n int) error {
	s.apply()

	if s.version == 0 {
		// init
		if n > 0 {
			// only up
			for i := 0; i < n; i++ {
				if err := s.do(i); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// gets current index of migrations
	idx := s.data.index(Migration{version: s.version})

	if n > 0 {
		// up
		for i := 0; i < n; i++ {
			if err := s.do(idx + i + 1); err != nil {
				return err
			}
		}
	} else {
		// down
		for i := 0; i < -n; i++ {
			if err := s.do(idx - i); err != nil {
				return err
			}
		}
	}

	return nil
}

// Apply applies the given migration version.
func (s *Service) Apply(d int, version uint64) error {
	s.direction = d
	s.apply()

	// gets current index of migrations
	idx := s.data.index(Migration{version: version})
	return s.do(idx)
}

func (s *Service) up(n int) error {
	s.direction = direction.Up
	err := s.step(n)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

func (s *Service) down(n int) error {
	s.direction = direction.Down
	err := s.step(-n)
	switch err {
	case errOutOfBoundsMigrations:
		return nil
	}
	return err
}

// Up upgrades migration version.
func (s *Service) Up() error {
	return s.up(notFoundIndex)
}

// Down downgrades migration version.
func (s *Service) Down() error {
	return s.down(notFoundIndex)
}

// Next upgrades migration version.
func (s *Service) Next(n int) error {
	return s.up(n)
}

// Prev downgrades migration version.
func (s *Service) Prev(n int) error {
	return s.down(n)
}

// Sync applies all migrations.
func (s *Service) Sync() error {
	s.direction = direction.Up
	s.apply()

	version := s.driver.Version()

	for _, mig := range s.data {
		if count := version.Count(mig.version); count == 0 {
			// gets current index of migrations
			idx := s.data.index(*mig)
			if err := s.do(idx); err != nil {
				return err
			}
		}
	}

	return nil
}

// NextMigration returns next version migrations.
func (s *Service) NextMigration(name string) (up *Migration, down *Migration, err error) {
	s.apply()

	// initialize default variables for making migrations.
	up, down = &Migration{version: 1, name: ""}, &Migration{version: 1, name: ""}
	verFormat := "%03d"

	// gets the oldest migration version file.
	if obj := s.data.last(); obj != nil {
		// for version number
		v := obj.version + 1
		up.version, down.version = v, v
	}
	if obj := s.data.first(); obj != nil {
		// for version format
		_, file := filepath.Split(obj.name)
		verFormat = version.Format(file)

		if _, err := time.Parse("20060102150405", version.Get(file)); err == nil {
			v := cast.Uint64(time.Now())
			up.version, down.version = v, v
		}

	}

	// [ver]_[name]_[direction-suffix][.ext]
	base := fmt.Sprintf("%s_%s_%%s%%s", verFormat, name)
	// including dot
	ext := s.driver.Ext()

	// up
	n := fmt.Sprintf(base, up.version, direction.Suffix(direction.Up), ext)
	up.name = filepath.Join(s.config.migrationsDir(), n)
	// down
	n = fmt.Sprintf(base, down.version, direction.Suffix(direction.Down), ext)
	down.name = filepath.Join(s.config.migrationsDir(), n)

	return
}
