package driver

import (
	"github.com/Fs02/kamimai/core"
)

const versionTableName = "schema_version"

func init() {
	core.RegisterDriver("mysql", &MySQL{})
}
