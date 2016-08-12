package driver

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/eure/kamimai/core"
	"github.com/lib/pq"
)

type (
	// Postgres driver object.
	Postgres struct {
		db *sql.DB
		tx *sql.Tx
		mu sync.Mutex
	}
)

// Open is the first function to be called.
// Check the dsn string and open and verify any connection
// that has to be made.
func (d *Postgres) Open(dsn string) error {
	z := strings.SplitN(dsn, "postgres:", 2)
	if len(z) != 2 {
		return errors.New("invalid data source name of postgres")
	}

	db, err := sql.Open("postgres", z[1])
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	d.db = db

	return d.Version().Create()
}

// Close is the last function to be called.
// Close any open connection here.
func (d *Postgres) Close() error {
	return d.db.Close()
}

// Ext returns the sql file extension used by path. The extension is the
// suffix beginning at the final dot in the final element of path; it is
// empty if there is no dot.
func (d *Postgres) Ext() string {
	return ".sql"
}

// Transaction starts a db transaction. The isolation level is dependent on the
// driver.
func (d *Postgres) Transaction(fn func(*sql.Tx) error) error {
	d.mu.Lock()
	defer func() {
		d.tx = nil
		d.mu.Unlock()
	}()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	d.tx = tx

	// Procedure
	if err := fn(d.tx); err != nil {
		if rberr := d.tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}

	// Commit
	if err := d.tx.Commit(); err != nil {
		if rberr := d.tx.Rollback(); rberr != nil {
			return rberr
		}
		return err
	}

	return nil
}

// Exec executes a query without returning any rows. The args are for any
// placeholder parameters in the query.
func (d *Postgres) Exec(query string, args ...interface{}) (sql.Result, error) {
	if d.tx != nil {
		return d.tx.Exec(query, args...)
	}
	return d.db.Exec(query, args...)
}

// Version returns a version interface.
func (d *Postgres) Version() core.Version {
	return d
}

// Migrate applies migration file.
func (d *Postgres) Migrate(m *core.Migration) error {
	b, err := m.Read()
	if err != nil {
		return err
	}

	stmts := bytes.Split(b, []byte(";"))
	for _, stmt := range stmts {
		query := strings.TrimSpace(string(stmt))
		if len(query) == 0 {
			continue
		}
		_, err = d.Exec(query)
		if err != nil {
			isWarn := strings.Contains(err.Error(), pq.Ewarning)
			if !isWarn {
				return err
			}
		}
	}

	return nil
}

// Insert inserts the given migration version.
func (d *Postgres) Insert(val uint64) error {
	query := fmt.Sprintf(`INSERT INTO %s (version) VALUES (%d)`,
		versionTableName, val)

	_, err := d.Exec(query)
	if err != nil {
		isWarn := strings.Contains(err.Error(), pq.Ewarning)
		if !isWarn {
			return err
		}
	}
	return nil
}

// Delete deletes the given migration version.
func (d *Postgres) Delete(val uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = %d`,
		versionTableName, val)

	_, err := d.Exec(query)
	if err != nil {
		isWarn := strings.Contains(err.Error(), pq.Ewarning)
		if !isWarn {
			return err
		}
	}
	return nil
}

// Count counts number of row the given migration version.
func (d *Postgres) Count(val uint64) int {
	query := fmt.Sprintf(`SELECT count(version) count FROM %s WHERE version = %d`,
		versionTableName, val)

	var count int
	if err := d.db.QueryRow(query).Scan(&count); err != nil {
		return 0
	}
	return count
}

// Current returns the current migration version.
func (d *Postgres) Current() (uint64, error) {
	const query = `SELECT version FROM ` +
		versionTableName + ` ORDER BY version DESC LIMIT 1`

	var version uint64
	err := d.db.QueryRow(query).Scan(&version)
	switch {
	case err == sql.ErrNoRows:
		return 0, nil
	case err != nil:
		return 0, err
	}
	return version, nil
}

// Create creates
func (d *Postgres) Create() error {
	const query = `CREATE TABLE IF NOT EXISTS ` +
		versionTableName + ` (version BIGINT NOT NULL PRIMARY KEY);`

	_, err := d.Exec(query)

	if err != nil {
		isWarn := strings.Contains(err.Error(), pq.Ewarning)
		if !isWarn {
			return err
		}
	}
	return nil
}

// Drop drops
func (d *Postgres) Drop() error {
	const query = `DROP TABLE IF EXISTS ` + versionTableName

	_, err := d.Exec(query)
	if err != nil {
		isWarn := strings.Contains(err.Error(), pq.Ewarning)
		if !isWarn {
			return err
		}
	}
	return nil
}
