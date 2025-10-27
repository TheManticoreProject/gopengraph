package properties_test

import (
	"reflect"
	"testing"

	"github.com/TheManticoreProject/gopengraph/properties"
)

func TestNewProperties(t *testing.T) {
	// Test empty constructor
	p1 := properties.NewProperties()
	if p1 == nil {
		t.Error("NewProperties() should not return nil")
	}
	if p1.Len() != 0 {
		t.Error("Empty Properties should have length 0")
	}

	// Test constructor with key-value pairs
	p2 := properties.NewProperties("name", "test", "age", 25, "active", true)
	if p2.Len() != 3 {
		t.Errorf("Expected length 3, got %d", p2.Len())
	}
	if p2.GetProperty("name") != "test" {
		t.Error("Expected name to be 'test'")
	}
	if p2.GetProperty("age") != 25 {
		t.Error("Expected age to be 25")
	}
	if p2.GetProperty("active") != true {
		t.Error("Expected active to be true")
	}

	// Test constructor with odd number of arguments (should ignore last one)
	p3 := properties.NewProperties("key1", "value1", "key2", "value2", "orphan")
	if p3.Len() != 2 {
		t.Errorf("Expected length 2, got %d", p3.Len())
	}
}

func TestSetProperty(t *testing.T) {
	p := properties.NewProperties()

	// Test valid property types
	validTests := []struct {
		key   string
		value interface{}
	}{
		{"string", "hello"},
		{"int", 42},
		{"float", 3.14},
		{"bool", true},
		{"nil", nil},
		{"slice", []string{"a", "b", "c"}},
	}

	for _, test := range validTests {
		p.SetProperty(test.key, test.value)
		if !p.HasProperty(test.key) {
			t.Errorf("Property %s should exist after setting", test.key)
		}
		got := p.GetProperty(test.key)
		if test.value != nil && reflect.TypeOf(test.value).Kind() == reflect.Slice {
			if !reflect.DeepEqual(got, test.value) {
				t.Errorf("Property %s should have value %v", test.key, test.value)
			}
		} else {
			if got != test.value {
				t.Errorf("Property %s should have value %v", test.key, test.value)
			}
		}
	}

	// Test invalid property types (should panic)
	invalidTests := []interface{}{
		map[string]string{"key": "value"}, // map
		struct{}{},                        // struct
		func() {},                         // function
	}

	for _, invalidValue := range invalidTests {
		func() {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Expected panic for invalid value %T, but no panic occurred", invalidValue)
				}
			}()
			p.SetProperty("invalid", invalidValue)
		}()
	}
}

func TestGetProperty(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)

	// Test getting existing property
	if p.GetProperty("name") != "test" {
		t.Error("Expected name to be 'test'")
	}

	// Test getting non-existing property without default
	if p.GetProperty("nonexistent") != nil {
		t.Error("Expected nonexistent property to return nil")
	}

	// Test getting non-existing property with default
	if p.GetProperty("nonexistent", "default") != "default" {
		t.Error("Expected nonexistent property to return default value")
	}

	// Test getting non-existing property with multiple defaults (should use first)
	if p.GetProperty("nonexistent", "first", "second") != "first" {
		t.Error("Expected nonexistent property to return first default value")
	}
}

func TestRemoveProperty(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)

	// Test removing existing property
	p.RemoveProperty("name")
	if p.HasProperty("name") {
		t.Error("Property 'name' should not exist after removal")
	}
	if p.Len() != 1 {
		t.Errorf("Expected length 1 after removal, got %d", p.Len())
	}

	// Test removing non-existing property (should not panic)
	p.RemoveProperty("nonexistent")
	if p.Len() != 1 {
		t.Errorf("Expected length to remain 1, got %d", p.Len())
	}
}

func TestHasProperty(t *testing.T) {
	p := properties.NewProperties("name", "test")

	// Test existing property
	if !p.HasProperty("name") {
		t.Error("Expected 'name' property to exist")
	}

	// Test non-existing property
	if p.HasProperty("nonexistent") {
		t.Error("Expected 'nonexistent' property to not exist")
	}
}

func TestGetAllProperties(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)
	props := p.GetAllProperties()

	// Test that we get a copy
	if len(props) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(props))
	}

	// Test that modifying the copy doesn't affect original
	props["new"] = "value"
	if p.HasProperty("new") {
		t.Error("Modifying returned properties should not affect original")
	}

	// Test that modifying original doesn't affect copy
	p.SetProperty("another", "value")
	if props["another"] != nil {
		t.Error("Modifying original should not affect returned copy")
	}
}

func TestClear(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)

	if p.Len() != 2 {
		t.Errorf("Expected length 2 before clear, got %d", p.Len())
	}

	p.Clear()

	if p.Len() != 0 {
		t.Errorf("Expected length 0 after clear, got %d", p.Len())
	}
	if p.HasProperty("name") {
		t.Error("Property 'name' should not exist after clear")
	}
}

func TestToDict(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)
	dict := p.ToDict()

	// Test that ToDict returns the same as GetAllProperties
	allProps := p.GetAllProperties()
	if len(dict) != len(allProps) {
		t.Error("ToDict should return same length as GetAllProperties")
	}

	for key, value := range dict {
		if allProps[key] != value {
			t.Errorf("ToDict value for key %s should match GetAllProperties", key)
		}
	}
}

func TestLen(t *testing.T) {
	p := properties.NewProperties()
	if p.Len() != 0 {
		t.Error("Empty Properties should have length 0")
	}

	p.SetProperty("key1", "value1")
	if p.Len() != 1 {
		t.Error("Properties with 1 item should have length 1")
	}

	p.SetProperty("key2", "value2")
	if p.Len() != 2 {
		t.Error("Properties with 2 items should have length 2")
	}

	p.RemoveProperty("key1")
	if p.Len() != 1 {
		t.Error("Properties should have length 1 after removing 1 item")
	}
}

func TestContains(t *testing.T) {
	p := properties.NewProperties("name", "test")

	// Test existing property
	if !p.Contains("name") {
		t.Error("Expected Contains to return true for existing property")
	}

	// Test non-existing property
	if p.Contains("nonexistent") {
		t.Error("Expected Contains to return false for non-existing property")
	}
}

func TestString(t *testing.T) {
	p := properties.NewProperties("name", "test", "age", 25)
	str := p.String()

	// Test that string representation contains expected elements
	expected := "Properties("
	if str[:len(expected)] != expected {
		t.Errorf("Expected string to start with %s, got %s", expected, str)
	}

	// Test that string representation contains the properties
	if len(str) <= len(expected) {
		t.Error("String representation should contain property information")
	}
}

func TestIsValidPropertyValue(t *testing.T) {
	p := properties.NewProperties()

	// Test valid types
	validTypes := []interface{}{
		"string",
		42,
		3.14,
		true,
		nil,
		[]string{"a", "b"},
		[]int{1, 2, 3},
	}

	for _, value := range validTypes {
		if !p.IsPropertyValueValid(value) {
			t.Errorf("Value %v of type %T should be valid", value, value)
		}
	}

	// Test invalid types
	invalidTypes := []interface{}{
		map[string]string{"key": "value"},
		struct{ field string }{field: "value"},
		func() {},
		make(chan int),
	}

	for _, value := range invalidTypes {
		if p.IsPropertyValueValid(value) {
			t.Errorf("Value %v of type %T should be invalid", value, value)
		}
	}
}

// Benchmark tests
func BenchmarkSetProperty(b *testing.B) {
	p := properties.NewProperties()
	for i := 0; i < b.N; i++ {
		p.SetProperty("key", "value")
	}
}

func BenchmarkGetProperty(b *testing.B) {
	p := properties.NewProperties("key", "value")
	for i := 0; i < b.N; i++ {
		p.GetProperty("key")
	}
}

func BenchmarkHasProperty(b *testing.B) {
	p := properties.NewProperties("key", "value")
	for i := 0; i < b.N; i++ {
		p.HasProperty("key")
	}
}
