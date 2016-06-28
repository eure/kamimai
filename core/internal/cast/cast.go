package cast

import (
	"strconv"
	"time"
)

// Uint64 casts an interface value to a uint64.
func Uint64(v interface{}) uint64 {
	switch v := v.(type) {
	case int:
		return (uint64)(v)
	case int8:
		return (uint64)(v)
	case int16:
		return (uint64)(v)
	case int32:
		return (uint64)(v)
	case int64:
		return (uint64)(v)
	case uint8:
		return (uint64)(v)
	case uint16:
		return (uint64)(v)
	case uint32:
		return (uint64)(v)
	case uint64:
		return v
	case string:
		n, _ := strconv.ParseUint(v, 10, 64)
		return n
	case *string:
		if v != nil {
			n, _ := strconv.ParseUint(*v, 10, 64)
			return n
		}
	case time.Time:
		if v.IsZero() {
			return 0
		}
		return Uint64(v.Format("20060102150405"))
	}
	return 0
}
