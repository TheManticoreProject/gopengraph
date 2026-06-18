package properties

import (
	"fmt"
	"reflect"
)

type Properties struct {
	Properties map[string]interface{}
}

// NewProperties creates a new Properties instance
func NewProperties() *Properties {
	p := &Properties{
		Properties: make(map[string]interface{}),
	}

	return p
}

// NewPropertiesFromMap creates a new Properties instance from a map of key-value pairs
func NewPropertiesFromMap(values map[string]interface{}) *Properties {
	p := NewProperties()

	for key, value := range values {
		p.SetProperty(key, value)
	}

	return p
}

func (p *Properties) SetProperty(key string, value interface{}) {
	if p.IsPropertyValueValid(value) {
		p.Properties[key] = value
	} else {
		panic(fmt.Sprintf("Property value must be a primitive type (string, int, float64, bool, nil, slice), got %T", value))
	}
}

func (p *Properties) GetProperty(key string, defaultVal ...interface{}) interface{} {
	if value, exists := p.Properties[key]; exists {
		return value
	}
	if len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return nil
}

func (p *Properties) RemoveProperty(key string) {
	delete(p.Properties, key)
}

func (p *Properties) HasProperty(key string) bool {
	_, exists := p.Properties[key]
	return exists
}

func (p *Properties) GetAllProperties() map[string]interface{} {
	// Return a copy to prevent external modification
	result := make(map[string]interface{})
	for k, v := range p.Properties {
		result[k] = v
	}
	return result
}

// Clear removes all properties
func (p *Properties) Clear() {
	p.Properties = make(map[string]interface{})
}

// IsPropertyValueValid reports whether value is a valid OpenGraph property value.
//
// The BloodHound OpenGraph schema restricts a property value to a single
// primitive (string, number, or boolean) or a homogeneous array of primitives.
// null values, nested objects, arrays of objects, and arrays mixing primitive
// types are not valid.
//
// Source: https://bloodhound.specterops.io/opengraph/developer/nodes
func (p *Properties) IsPropertyValueValid(value interface{}) bool {
	if value == nil {
		return false
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice, reflect.Array:
		return isHomogeneousPrimitiveSequence(reflect.ValueOf(value))
	default:
		return primitiveCategory(value) != ""
	}
}

// primitiveCategory classifies value into one of the OpenGraph primitive
// categories ("string", "number", "boolean"). It returns "" when value is nil
// or not a primitive (e.g. an object, slice, or array).
func primitiveCategory(value interface{}) string {
	if value == nil {
		return ""
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	default:
		return ""
	}
}

// isHomogeneousPrimitiveSequence reports whether v is a slice or array whose
// elements are all primitives belonging to the same OpenGraph category. An
// empty sequence is considered valid.
func isHomogeneousPrimitiveSequence(v reflect.Value) bool {
	category := ""
	for i := 0; i < v.Len(); i++ {
		c := primitiveCategory(v.Index(i).Interface())
		if c == "" {
			return false
		}
		if category == "" {
			category = c
		} else if c != category {
			return false
		}
	}
	return true
}

// ToDict converts properties to map for JSON serialization
func (p *Properties) ToDict() map[string]interface{} {
	return p.GetAllProperties()
}

// Len returns the number of properties
func (p *Properties) Len() int {
	return len(p.Properties)
}

// Contains checks if a key exists
func (p *Properties) Contains(key string) bool {
	return p.HasProperty(key)
}

// String returns string representation of Properties
func (p *Properties) String() string {
	return fmt.Sprintf("Properties(%v)", p.Properties)
}
