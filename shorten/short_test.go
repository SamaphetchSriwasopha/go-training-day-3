package shorten

import (
	"testing"
)

func TestNewShorter(t *testing.T) {
	shorter := NewShorter()

	// Generate a shorten code
	code1, err := shorter.generate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if code1 == "" {
		t.Error("expected non-empty string code")
	}

	// Generate a second code and check for uniqueness (using rand should be unique)
	code2, err := shorter.generate()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if code1 == code2 {
		t.Errorf("expected unique codes, but got identical codes: %s", code1)
	}
}
