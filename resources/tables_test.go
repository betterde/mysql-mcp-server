package resources

import (
	"testing"
)

func TestExtractDatabase(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    string
		wantErr bool
	}{
		{"simple database", "mysql://mydb/tables", "mydb", false},
		{"nested path", "mysql://testdb/tables/extra", "testdb", false},
		{"empty database", "mysql:///tables", "", true},
		{"missing database", "mysql://", "", true},
		{"invalid uri", "://", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractDatabase(tt.uri)
			if tt.wantErr {
				if err == nil {
					t.Errorf("extractDatabase(%q) = %q, want error", tt.uri, got)
				}
				return
			}
			if err != nil {
				t.Errorf("extractDatabase(%q) unexpected error: %v", tt.uri, err)
				return
			}
			if got != tt.want {
				t.Errorf("extractDatabase(%q) = %q, want %q", tt.uri, got, tt.want)
			}
		})
	}
}

func TestExtractDatabaseAndTable(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantDB  string
		wantTbl string
		wantErr bool
	}{
		{"columns uri", "mysql://mydb/users/columns", "mydb", "users", false},
		{"ddl uri", "mysql://mydb/users/ddl", "mydb", "users", false},
		{"extra path segments", "mysql://mydb/users/columns/extra", "mydb", "users", false},
		{"missing database", "mysql:///users/columns", "", "", true},
		{"only one segment", "mysql://mydb", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, tbl, err := extractDatabaseAndTable(tt.uri)
			if tt.wantErr {
				if err == nil {
					t.Errorf("extractDatabaseAndTable(%q) = (%q, %q), want error", tt.uri, db, tbl)
				}
				return
			}
			if err != nil {
				t.Errorf("extractDatabaseAndTable(%q) unexpected error: %v", tt.uri, err)
				return
			}
			if db != tt.wantDB {
				t.Errorf("extractDatabaseAndTable(%q) db = %q, want %q", tt.uri, db, tt.wantDB)
			}
			if tbl != tt.wantTbl {
				t.Errorf("extractDatabaseAndTable(%q) tbl = %q, want %q", tt.uri, tbl, tt.wantTbl)
			}
		})
	}
}
