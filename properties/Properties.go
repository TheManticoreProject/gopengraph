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

// isValidPropertyValue checks if a value is a valid property type
func (p *Properties) IsPropertyValueValid(value interface{}) bool {
	if value == nil {
		return true
	}

	valueType := reflect.TypeOf(value)
	switch valueType.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.Slice:
		return true
	default:
		return false
	}
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
