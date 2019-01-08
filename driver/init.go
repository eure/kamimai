package driver

import (
	"github.com/mr04vv/kamimai/core"
)

const versionTableName = "schema_version"

func init() {
	core.RegisterDriver("mysql", &MySQL{})
}
