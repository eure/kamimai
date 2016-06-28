package driver

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/kaneshin/kamimai/core"
)

type (
	// MySQL driver object.
	MySQL struct {
		db *sql.DB
	}
)

// Open is the first function to be called.
// Check the dsn string and open and verify any connection
// that has to be made.
func (d *MySQL) Open(dsn string) error {
	z := strings.SplitN(dsn, "mysql://", 2)
	if len(z) != 2 {
		return errors.New("invalid data source name of mysql")
	}

	db, err := sql.Open("mysql", z[1])
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
func (d *MySQL) Close() error {
	return d.db.Close()
}

// Ext returns the sql file extension used by path. The extension is the
// suffix beginning at the final dot in the final element of path; it is
// empty if there is no dot.
func (d *MySQL) Ext() string {
	return ".sql"
}

// Version returns a version interface.
func (d *MySQL) Version() core.Version {
	return d
}

// Migrate applies migration file.
func (d *MySQL) Migrate(m *core.Migration) error {
	b, err := m.Read()
	if err != nil {
		return err
	}
	_, err = d.db.Exec(string(b))
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}
	if err := d.Insert(m.Version()); err != nil {
		return err
	}
	return nil
}

// Insert inserts the given migration version.
func (d *MySQL) Insert(val uint64) error {
	query := fmt.Sprintf(`INSERT INTO %s (version) VALUES (%d)`,
		versionTableName, val)

	_, err := d.db.Exec(query)
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}
	return nil
}

// Delete deletes the given migration version.
func (d *MySQL) Delete(val uint64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE version = %d`,
		versionTableName, val)

	_, err := d.db.Exec(query)
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}
	return nil
}

// Current returns the current migration version.
func (d *MySQL) Current() (uint64, error) {
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
func (d *MySQL) Create() error {
	const query = `CREATE TABLE IF NOT EXISTS ` +
		versionTableName + ` (version int not null primary key);`

	_, err := d.db.Exec(query)
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}
	return nil
}

// Drop drops
func (d *MySQL) Drop() error {
	const query = `DROP TABLE IF EXISTS ` + versionTableName

	_, err := d.db.Exec(query)
	if _, isWarn := err.(mysql.MySQLWarnings); err != nil && !isWarn {
		return err
	}
	return nil
}
