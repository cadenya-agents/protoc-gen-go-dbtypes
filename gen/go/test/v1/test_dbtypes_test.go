package testv1

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestToolSetSpecValue_RoundTrip(t *testing.T) {
	// Create a test message
	spec := &ToolSetSpec{
		ToolIds: []string{"tool-1", "tool-2", "tool-3"},
		Name:    "my-toolset",
		Enabled: true,
	}

	// Wrap it
	wrapper := NewToolSetSpecValue(spec)

	// Get database value
	dbVal, err := wrapper.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}

	// Scan back into a new wrapper
	wrapper2 := &ToolSetSpecValue{}
	if err := wrapper2.Scan(dbVal); err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	// Verify
	if !proto.Equal(spec, wrapper2.Unwrap()) {
		t.Errorf("round-trip failed:\ngot:  %v\nwant: %v", wrapper2.Unwrap(), spec)
	}
}

func TestToolSetSpecValue_NilMessage(t *testing.T) {
	// Test with nil message
	wrapper := NewToolSetSpecValue(nil)

	// Should still be able to get value
	val, err := wrapper.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}

	// Empty proto should produce some bytes
	if val == nil {
		t.Error("Value() returned nil for empty message")
	}
}

func TestToolSetSpecValue_NilWrapper(t *testing.T) {
	// Test with nil ProtoValue
	wrapper := &ToolSetSpecValue{}

	// Value should return nil for nil wrapper
	val, err := wrapper.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}
	if val != nil {
		t.Errorf("Value() = %v, want nil", val)
	}
}

func TestToolSetSpecValue_ScanNil(t *testing.T) {
	wrapper := &ToolSetSpecValue{}

	// Scanning nil should not error
	if err := wrapper.Scan(nil); err != nil {
		t.Fatalf("Scan(nil) error: %v", err)
	}
}

func TestToolSetSpecValue_ScanString(t *testing.T) {
	spec := &ToolSetSpec{
		ToolIds: []string{"tool-1"},
		Name:    "test",
		Enabled: true,
	}

	// Marshal to bytes
	data, err := proto.Marshal(spec)
	if err != nil {
		t.Fatalf("proto.Marshal error: %v", err)
	}

	// Scan from string (some databases return strings)
	wrapper := &ToolSetSpecValue{}
	if err := wrapper.Scan(string(data)); err != nil {
		t.Fatalf("Scan(string) error: %v", err)
	}

	if !proto.Equal(spec, wrapper.Unwrap()) {
		t.Errorf("scan from string failed:\ngot:  %v\nwant: %v", wrapper.Unwrap(), spec)
	}
}

func TestToolSetSpecValue_ScanInvalidType(t *testing.T) {
	wrapper := &ToolSetSpecValue{}

	err := wrapper.Scan(123) // invalid type
	if err == nil {
		t.Error("Scan(int) should return error")
	}
}

func TestToolSetSpecValue_Unwrap(t *testing.T) {
	spec := &ToolSetSpec{
		Name: "test",
	}

	wrapper := NewToolSetSpecValue(spec)
	unwrapped := wrapper.Unwrap()

	if unwrapped != spec {
		t.Error("Unwrap() should return the same message pointer")
	}
}

func TestToolSetSpecValue_UnwrapNil(t *testing.T) {
	wrapper := &ToolSetSpecValue{}
	unwrapped := wrapper.Unwrap()

	if unwrapped != nil {
		t.Errorf("Unwrap() = %v, want nil", unwrapped)
	}
}

func TestUserPreferencesValue_RoundTrip(t *testing.T) {
	prefs := &UserPreferences{
		Theme:    "dark",
		Language: "en",
		Settings: map[string]string{
			"notifications": "enabled",
			"autoSave":      "true",
		},
	}

	wrapper := NewUserPreferencesValue(prefs)

	dbVal, err := wrapper.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}

	wrapper2 := &UserPreferencesValue{}
	if err := wrapper2.Scan(dbVal); err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !proto.Equal(prefs, wrapper2.Unwrap()) {
		t.Errorf("round-trip failed:\ngot:  %v\nwant: %v", wrapper2.Unwrap(), prefs)
	}
}

func TestContainerValue_RoundTrip(t *testing.T) {
	container := &Container{
		Id: "container-1",
		Spec: &ToolSetSpec{
			ToolIds: []string{"tool-a"},
			Name:    "nested-spec",
		},
		Items: []*Container_Item{
			{Key: "key1", Value: "value1"},
			{Key: "key2", Value: "value2"},
		},
	}

	wrapper := NewContainerValue(container)

	dbVal, err := wrapper.Value()
	if err != nil {
		t.Fatalf("Value() error: %v", err)
	}

	wrapper2 := &ContainerValue{}
	if err := wrapper2.Scan(dbVal); err != nil {
		t.Fatalf("Scan() error: %v", err)
	}

	if !proto.Equal(container, wrapper2.Unwrap()) {
		t.Errorf("round-trip failed:\ngot:  %v\nwant: %v", wrapper2.Unwrap(), container)
	}
}
