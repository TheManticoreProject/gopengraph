package edge_test

import (
	"testing"

	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/properties"
)

func TestNewEdge(t *testing.T) {
	tests := []struct {
		name          string
		startNodeID   string
		endNodeID     string
		kind          string
		properties    *properties.Properties
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid edge",
			startNodeID: "start1",
			endNodeID:   "end1",
			kind:        "CONNECTS_TO",
			properties:  properties.NewProperties(),
			expectError: false,
		},
		{
			name:          "empty start node ID",
			startNodeID:   "",
			endNodeID:     "end1",
			kind:          "CONNECTS_TO",
			expectError:   true,
			errorContains: "start node ID cannot be empty",
		},
		{
			name:          "empty end node ID",
			startNodeID:   "start1",
			endNodeID:     "",
			kind:          "CONNECTS_TO",
			expectError:   true,
			errorContains: "end node ID cannot be empty",
		},
		{
			name:          "empty kind",
			startNodeID:   "start1",
			endNodeID:     "end1",
			kind:          "",
			expectError:   true,
			errorContains: "edge kind cannot be empty",
		},
		{
			name:        "nil properties",
			startNodeID: "start1",
			endNodeID:   "end1",
			kind:        "CONNECTS_TO",
			properties:  nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := edge.NewEdge(tt.startNodeID, tt.endNodeID, tt.kind, tt.properties)

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

			if e.GetStartNodeID() != tt.startNodeID {
				t.Errorf("expected start node ID %q, got %q", tt.startNodeID, e.GetStartNodeID())
			}
			if e.GetEndNodeID() != tt.endNodeID {
				t.Errorf("expected end node ID %q, got %q", tt.endNodeID, e.GetEndNodeID())
			}
			if e.GetKind() != tt.kind {
				t.Errorf("expected kind %q, got %q", tt.kind, e.GetKind())
			}
		})
	}
}

func TestEdgeProperties(t *testing.T) {
	e, err := edge.NewEdge("start1", "end1", "CONNECTS_TO", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test setting and getting properties
	e.SetProperty("weight", 10)
	if val := e.GetProperty("weight"); val != 10 {
		t.Errorf("expected property value 10, got %v", val)
	}

	// Test getting non-existent property with default
	if val := e.GetProperty("nonexistent", "default"); val != "default" {
		t.Errorf("expected default value 'default', got %v", val)
	}

	// Test removing property
	e.RemoveProperty("weight")
	if val := e.GetProperty("weight"); val != nil {
		t.Errorf("expected nil after removal, got %v", val)
	}
}

func TestEdgeEqual(t *testing.T) {
	e1, _ := edge.NewEdge("start1", "end1", "CONNECTS_TO", nil)
	e2, _ := edge.NewEdge("start1", "end1", "CONNECTS_TO", nil)
	e3, _ := edge.NewEdge("start2", "end2", "CONNECTS_TO", nil)

	if !e1.Equal(e2) {
		t.Error("expected identical edges to be equal")
	}

	if e1.Equal(e3) {
		t.Error("expected different edges to not be equal")
	}

	if e1.Equal(nil) {
		t.Error("expected edge compared with nil to not be equal")
	}
}

func TestEdgeToDict(t *testing.T) {
	props := properties.NewProperties()
	props.SetProperty("weight", 10)

	e, _ := edge.NewEdge("start1", "end1", "CONNECTS_TO", props)

	dict := e.ToDict()

	// Check basic structure
	if dict["kind"] != "CONNECTS_TO" {
		t.Errorf("expected kind 'CONNECTS_TO', got %v", dict["kind"])
	}

	// Check start node
	startNode, ok := dict["start"].(map[string]string)
	if !ok {
		t.Fatal("expected start to be map[string]string")
	}
	if startNode["value"] != "start1" || startNode["match_by"] != "id" {
		t.Error("incorrect start node mapping")
	}

	// Check end node
	endNode, ok := dict["end"].(map[string]string)
	if !ok {
		t.Fatal("expected end to be map[string]string")
	}
	if endNode["value"] != "end1" || endNode["match_by"] != "id" {
		t.Error("incorrect end node mapping")
	}

	// Check properties
	propsDict, ok := dict["properties"].(map[string]interface{})
	if !ok {
		t.Fatal("expected properties to be map[string]interface{}")
	}
	if propsDict["weight"] != 10 {
		t.Errorf("expected weight 10, got %v", propsDict["weight"])
	}
}

func TestEdgeString(t *testing.T) {
	e, _ := edge.NewEdge("start1", "end1", "CONNECTS_TO", nil)
	str := e.String()
	expected := "Edge(start='start1', end='end1', kind='CONNECTS_TO', properties=map[])"
	if str != expected {
		t.Errorf("expected string %q, got %q", expected, str)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-1] != substr[0]
}
