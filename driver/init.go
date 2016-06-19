package driver

import (
	"github.com/kaneshin/kamimai/core"
)

const versionTableName = "schema_version"

func init() {
	core.RegisterDriver("mysql", &MySQL{})
}
