package query

import (
	"testing"
)

// Test that write queries are rejected by normalizeReadOnlyQuery (the read-only guard)
// while they would be accepted by the execute tool's own validation.
func TestExecuteVsSelectDistinction(t *testing.T) {
	writeQueries := []string{
		"INSERT INTO t VALUES (1)",
		"UPDATE t SET c=1",
		"DELETE FROM t",
		"CREATE TABLE t (id INT)",
		"ALTER TABLE t ADD COLUMN c INT",
		"DROP TABLE t",
		"TRUNCATE TABLE t",
		"RENAME TABLE t TO t2",
	}

	for _, q := range writeQueries {
		t.Run(q, func(t *testing.T) {
			_, err := normalizeReadOnlyQuery(q)
			if err == nil {
				t.Errorf("normalizeReadOnlyQuery(%q) should reject write queries", q)
			}
		})
	}
}

func TestNormalizeExplainFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  string
		want    string
		wantErr bool
	}{
		{"empty", "", "", false},
		{"traditional", "TRADITIONAL", "TRADITIONAL", false},
		{"tree", "TREE", "TREE", false},
		{"json", "JSON", "JSON", false},
		{"lowercase tree", "tree", "TREE", false},
		{"mixed case", "Json", "JSON", false},
		{"invalid format", "XML", "", true},
		{"generic", "graph", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeExplainFormat(tt.format)
			if tt.wantErr {
				if err == nil {
					t.Errorf("normalizeExplainFormat(%q) = %q, want error", tt.format, got)
				}
				return
			}
			if err != nil {
				t.Errorf("normalizeExplainFormat(%q) unexpected error: %v", tt.format, err)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeExplainFormat(%q) = %q, want %q", tt.format, got, tt.want)
			}
		})
	}
}
