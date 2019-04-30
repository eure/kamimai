package driver

import (
	"github.com/eure/kamimai/core"
)

var versionTableName = "schema_version"

func init() {
	core.RegisterDriver("mysql", &MySQL{})
}

func SetVersionTable(n string) {
	versionTableName = n
}
