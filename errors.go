package smooth

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrMigrationSystemDisabled = errors.New("migration system is disabled")
)
