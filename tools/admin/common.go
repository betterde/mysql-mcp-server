package admin

import "database/sql"

const (
	defaultAdminLimit = 100
	maxAdminLimit     = 1000
)

func normalizeAdminLimit(limit int) int {
	if limit <= 0 {
		return defaultAdminLimit
	}
	if limit > maxAdminLimit {
		return maxAdminLimit
	}
	return limit
}

func nullStringValue(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
