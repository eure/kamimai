package driver

import (
	"github.com/eure/kamimai/core"
)

const versionTableName = "schema_version"

func init() {
	core.RegisterDriver("mysql", &MySQL{})
	core.RegisterDriver("postgres", &Postgres{})
}
