package node_test

import (
	"testing"

	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

func TestNewNode(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		kinds         []string
		properties    *properties.Properties
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid node",
			id:          "node1",
			kinds:       []string{"User", "Group"},
			properties:  properties.NewProperties(),
			expectError: false,
		},
		{
			name:          "empty id",
			id:            "",
			kinds:         []string{"User"},
			expectError:   true,
			errorContains: "node ID cannot be empty",
		},
		{
			name:        "nil kinds",
			id:          "node1",
			kinds:       nil,
			properties:  properties.NewProperties(),
			expectError: false,
		},
		{
			name:        "nil properties",
			id:          "node1",
			kinds:       []string{"User"},
			properties:  nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := node.NewNode(tt.id, tt.kinds, tt.properties)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorContains)
				} else if err.Error() != tt.errorContains {
					t.Errorf("expected error containing %q, got %q", tt.errorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if n.GetID() != tt.id {
				t.Errorf("expected ID %q, got %q", tt.id, n.GetID())
			}

			if tt.kinds == nil {
				if len(n.GetKinds()) != 0 {
					t.Errorf("expected empty kinds slice, got %v", n.GetKinds())
				}
			} else {
				for _, kind := range tt.kinds {
					if !n.HasKind(kind) {
						t.Errorf("expected node to have kind %q", kind)
					}
				}
			}
		})
	}
}

func TestNodeKinds(t *testing.T) {
	n, err := node.NewNode("node1", []string{"User"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test adding new kind
	n.AddKind("Group")
	if !n.HasKind("Group") {
		t.Error("expected node to have kind 'Group'")
	}

	// Test adding duplicate kind
	n.AddKind("User")
	kinds := n.GetKinds()
	count := 0
	for _, k := range kinds {
		if k == "User" {
			count++
		}
	}
	if count > 1 {
		t.Error("expected no duplicate kinds")
	}

	// Test removing kind
	n.RemoveKind("User")
	if n.HasKind("User") {
		t.Error("expected node to not have kind 'User'")
	}
}

func TestNodeProperties(t *testing.T) {
	n, err := node.NewNode("node1", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test setting and getting properties
	n.SetProperty("name", "test")
	if val := n.GetProperty("name"); val != "test" {
		t.Errorf("expected property value 'test', got %v", val)
	}

	// Test getting non-existent property with default
	if val := n.GetProperty("nonexistent", "default"); val != "default" {
		t.Errorf("expected default value 'default', got %v", val)
	}

	// Test removing property
	n.RemoveProperty("name")
	if val := n.GetProperty("name"); val != nil {
		t.Errorf("expected nil value after removal, got %v", val)
	}
}

func TestNodeToDict(t *testing.T) {
	props := properties.NewProperties()
	props.SetProperty("name", "test")

	n, _ := node.NewNode("node1", []string{"User", "Group"}, props)

	dict := n.ToDict()

	if dict["id"] != "node1" {
		t.Errorf("expected id 'node1', got %v", dict["id"])
	}

	kinds, ok := dict["kinds"].([]string)
	if !ok {
		t.Fatal("expected kinds to be []string")
	}
	if len(kinds) != 2 || !contains(kinds, "User") || !contains(kinds, "Group") {
		t.Error("incorrect kinds in dictionary")
	}

	propsDict, ok := dict["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("expected properties to be map[string]interface{}")
	}
	if propsDict["name"] != "test" {
		t.Errorf("expected name 'test', got %v", propsDict["name"])
	}
}

func TestNodeEqual(t *testing.T) {
	n1, _ := node.NewNode("node1", nil, nil)
	n2, _ := node.NewNode("node1", []string{"User"}, nil)
	n3, _ := node.NewNode("node2", nil, nil)

	if !n1.Equal(n2) {
		t.Error("expected nodes with same ID to be equal")
	}

	if n1.Equal(n3) {
		t.Error("expected nodes with different IDs to not be equal")
	}

	if n1.Equal(nil) {
		t.Error("expected node compared with nil to not be equal")
	}
}

func TestNodeString(t *testing.T) {
	n, _ := node.NewNode("node1", []string{"User"}, nil)
	str := n.String()
	expected := "Node(id='node1', kinds=[User], properties=map[])"
	if str != expected {
		t.Errorf("expected string %q, got %q", expected, str)
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
