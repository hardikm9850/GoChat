package mysql

import "strings"

func isDuplicateKeyError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry")
}
