//go:build unit

package cmd

import (
	"testing"
)

func TestShow_JSON(t *testing.T) {
	data := map[string]string{"key": "value"}

	if err := show("json", data); err != nil {
		t.Fatalf("show(json) returned error: %v", err)
	}
}

func TestShow_YAML(t *testing.T) {
	data := map[string]string{"key": "value"}

	if err := show("yaml", data); err != nil {
		t.Fatalf("show(yaml) returned error: %v", err)
	}
}

func TestShow_DefaultIsJSON(t *testing.T) {
	data := map[string]string{"key": "value"}

	if err := show("unknown", data); err != nil {
		t.Fatalf("show(unknown) returned error: %v", err)
	}
}

func TestShow_JSONMarshalError(t *testing.T) {
	// channels cannot be marshaled to JSON
	data := make(chan int)

	if err := show("json", data); err == nil {
		t.Fatal("show(json) with unmarshalable data should return error")
	}
}

func TestShow_EmptyStruct(t *testing.T) {
	type empty struct{}

	if err := show("json", empty{}); err != nil {
		t.Fatalf("show(json) with empty struct returned error: %v", err)
	}

	if err := show("yaml", empty{}); err != nil {
		t.Fatalf("show(yaml) with empty struct returned error: %v", err)
	}
}
