package handler

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestResourceResult(t *testing.T) {
	tests := []struct {
		name string
		data any
	}{
		{"string map", map[string]string{"key": "value"}},
		{"int value", 42},
		{"slice", []string{"a", "b", "c"}},
		{"struct", struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}{Name: "test", Age: 30}},
		{"nil", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ResourceResult(tt.data)
			if err != nil {
				t.Fatalf("ResourceResult(%v) unexpected error: %v", tt.data, err)
			}
			if result == nil {
				t.Fatal("ResourceResult returned nil")
			}
			if len(result.Contents) != 1 {
				t.Fatalf("expected 1 content, got %d", len(result.Contents))
			}
			if result.Contents[0].MIMEType != "application/json" {
				t.Errorf("expected MIMEType 'application/json', got %q", result.Contents[0].MIMEType)
			}
			if result.Contents[0].Text == "" && tt.data != nil {
				t.Error("expected non-empty Text for non-nil data")
			}

			var decoded any
			if err := json.Unmarshal([]byte(result.Contents[0].Text), &decoded); err != nil {
				t.Errorf("Text is not valid JSON: %v\nText: %s", err, result.Contents[0].Text)
			}
		})
	}
}

func TestResourceResultMarshalError(t *testing.T) {
	ch := make(chan int)
	_, err := ResourceResult(ch)
	if err == nil {
		t.Error("expected error for unmarshalable type, got nil")
	}
	if !strings.Contains(err.Error(), "json") {
		t.Errorf("expected json marshal error, got %v", err)
	}
}
