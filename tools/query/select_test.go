package query

import (
	"testing"
)

func TestNormalizeReadOnlyQuery(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		want    string
		wantErr bool
	}{
		{"simple select", "SELECT 1", "SELECT 1", false},
		{"select with from", "SELECT * FROM users", "SELECT * FROM users", false},
		{"show databases", "SHOW DATABASES", "SHOW DATABASES", false},
		{"describe table", "DESCRIBE t", "DESCRIBE t", false},
		{"desc table", "DESC t", "DESC t", false},
		{"with cte", "WITH cte AS (SELECT 1) SELECT * FROM cte", "WITH cte AS (SELECT 1) SELECT * FROM cte", false},
		{"trailing whitespace", "  SELECT 1  ", "SELECT 1", false},
		{"empty query", "", "", true},
		{"insert rejected", "INSERT INTO t VALUES (1)", "", true},
		{"update rejected", "UPDATE t SET x=1", "", true},
		{"delete rejected", "DELETE FROM t", "", true},
		{"drop rejected", "DROP TABLE t", "", true},
		{"create rejected", "CREATE TABLE t (id INT)", "", true},
		{"alter rejected", "ALTER TABLE t ADD COLUMN c INT", "", true},
		{"multi statement", "SELECT 1; SELECT 2", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeReadOnlyQuery(tt.query)
			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeReadOnlyQuery(%q) = %q, want error", tt.query, got)
				}
				return
			}
			if err != nil {
				t.Errorf("normalizeReadOnlyQuery(%q) unexpected error: %v", tt.query, err)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeReadOnlyQuery(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestNormalizeLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"zero defaults", 0, 100},
		{"negative defaults", -1, 100},
		{"within range", 50, 50},
		{"at default", 100, 100},
		{"above max clamps", 2000, 1000},
		{"at max", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeLimit(tt.limit)
			if got != tt.want {
				t.Errorf("normalizeLimit(%d) = %d, want %d", tt.limit, got, tt.want)
			}
		})
	}
}

func TestFirstKeyword(t *testing.T) {
	tests := []struct {
		query string
		want  string
	}{
		{"SELECT * FROM t", "select"},
		{"  SHOW DATABASES", "show"},
		{"DESCRIBE t", "describe"},
		{"WITH cte AS (...)", "with"},
		{"INSERT INTO t VALUES(1)", "insert"},
		{"(SELECT 1)", "select"},
		{"", ""},
		{" \t\r\n", ""},
	}

	for _, tt := range tests {
		t.Run(tt.query, func(t *testing.T) {
			got := firstKeyword(tt.query)
			if got != tt.want {
				t.Errorf("firstKeyword(%q) = %q, want %q", tt.query, got, tt.want)
			}
		})
	}
}

func TestNormalizeValue(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  any
	}{
		{"nil", nil, nil},
		{"bytes to string", []byte("hello"), "hello"},
		{"int64 unchanged", int64(42), int64(42)},
		{"uint64 unchanged", uint64(42), uint64(42)},
		{"float64 unchanged", float64(3.14), float64(3.14)},
		{"bool unchanged", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeValue(tt.value)
			if got != tt.want {
				t.Errorf("normalizeValue(%v) = %v (type %T), want %v (type %T)", tt.value, got, got, tt.want, tt.want)
			}
		})
	}
}
