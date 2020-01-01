package core

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/eure/kamimai/internal/cast"
	"github.com/eure/kamimai/internal/direction"
	"github.com/eure/kamimai/internal/version"
)

const (
	notFoundIndex = 0xffff
)

var (
	errOutOfBoundsMigrations     = errors.New("out of bounds migration")
	errDuplicateMigrationVersion = errors.New("duplicate migration version")
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

func (s Service) walker(indexPath map[uint64]*Migration) func(string, os.FileInfo, error) error {
	wd, _ := os.Getwd()

	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if s.direction != direction.Get(name) {
			return nil
		}

		fullname := filepath.Clean(filepath.Join(wd, path))
		ver := cast.Uint64(version.Get(name))
		mig, found := indexPath[ver]
		if found && mig.IsValid() {
			return errors.Wrap(errDuplicateMigrationVersion, fmt.Sprintf("failed to read migration %s", fullname))
		}

		mig = &Migration{
			version: ver,
		}
		indexPath[ver] = mig
		mig.name = fullname

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

func (s *Service) mustApply() {
	if err := s.apply(); err != nil {
		panic(err)
	}
}

func (s *Service) apply() error {
	index := map[uint64]*Migration{}
	if err := filepath.Walk(s.config.migrationsDir(), s.walker(index)); err != nil {
		return err
	}

	list := make([]*Migration, len(index))
	i := 0
	for _, m := range index {
		list[i] = m
		i++
	}

	migs := Migrations(list)
	sort.Sort(migs)
	s.data = migs
	return nil
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
	if err := s.apply(); err != nil {
		return err
	}

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
	if err := s.apply(); err != nil {
		return err
	}

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
	if err := s.apply(); err != nil {
		return err
	}

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
	if err := s.apply(); err != nil {
		return nil, nil, err
	}

	// initialize default variables for making migrations.
	up, down = &Migration{version: 1, name: ""}, &Migration{version: 1, name: ""}
	verFormat := "%03d"

	// check if the format of the latest version is timestamp
	isDateFormat := IsTimeStamp(s.version)

	if isDateFormat {
		// for version format
		v := cast.Uint64(time.Now())
		up.version, down.version = v, v
	} else {
		// for version number
		v := s.version + 1
		up.version, down.version = v, v
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

// IsTimeStamp return if prefix is timestamp or not
func IsTimeStamp(prefix uint64) (result bool) {
	date := strconv.FormatUint(prefix, 10)
	layout := "20060102150405"
	t, _ := time.Parse(layout, date)
	std, _ := time.Parse(layout, layout)
	result = t.After(std)
	return
}
