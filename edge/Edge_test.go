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
	startNode, ok := dict["start"].(map[string]interface{})
	if !ok {
		t.Fatal("expected start to be map[string]interface{}")
	}
	if startNode["value"] != "start1" || startNode["match_by"] != "id" {
		t.Error("incorrect start node mapping")
	}

	// Check end node
	endNode, ok := dict["end"].(map[string]interface{})
	if !ok {
		t.Fatal("expected end to be map[string]interface{}")
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

func TestEdgeMatchByName(t *testing.T) {
	e, err := edge.NewEdgeWithEndpoints(
		edge.NewEndpointByName("alice", "User"),
		edge.NewEndpointByName("file-server-1", "Server"),
		"HasAccess",
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dict := e.ToDict()
	start, ok := dict["start"].(map[string]interface{})
	if !ok {
		t.Fatal("expected start to be map[string]interface{}")
	}
	if start["match_by"] != edge.MatchByName {
		t.Errorf("expected match_by %q, got %v", edge.MatchByName, start["match_by"])
	}
	if start["value"] != "alice" {
		t.Errorf("expected value 'alice', got %v", start["value"])
	}
	if start["kind"] != "User" {
		t.Errorf("expected kind 'User', got %v", start["kind"])
	}
	if _, present := start["property_matchers"]; present {
		t.Error("name endpoint must not include property_matchers")
	}
}

func TestEdgeMatchByProperty(t *testing.T) {
	matchers := []edge.PropertyMatcher{
		{Key: "username", Operator: "equals", Value: "alice.smith"},
		{Key: "active", Operator: "equals", Value: true},
	}
	e, err := edge.NewEdgeWithEndpoints(
		edge.NewEndpointByProperty(matchers, "User"),
		edge.NewEndpointByID("server-1"),
		"CustomRelationship",
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dict := e.ToDict()
	start, ok := dict["start"].(map[string]interface{})
	if !ok {
		t.Fatal("expected start to be map[string]interface{}")
	}
	if start["match_by"] != edge.MatchByProperty {
		t.Errorf("expected match_by %q, got %v", edge.MatchByProperty, start["match_by"])
	}
	if _, present := start["value"]; present {
		t.Error("property endpoint must not include a top-level value")
	}
	pm, ok := start["property_matchers"].([]map[string]interface{})
	if !ok {
		t.Fatalf("expected property_matchers to be []map[string]interface{}, got %T", start["property_matchers"])
	}
	if len(pm) != 2 {
		t.Fatalf("expected 2 property matchers, got %d", len(pm))
	}
	if pm[0]["key"] != "username" || pm[0]["operator"] != "equals" || pm[0]["value"] != "alice.smith" {
		t.Errorf("unexpected first matcher: %v", pm[0])
	}
}

func TestEndpointValidate(t *testing.T) {
	// id/name require a value
	if err := edge.NewEndpointByID("").Validate(); err == nil {
		t.Error("expected error for id endpoint with empty value")
	}
	// property requires at least one matcher
	if err := edge.NewEndpointByProperty(nil, "User").Validate(); err == nil {
		t.Error("expected error for property endpoint without matchers")
	}
	// property matcher requires a key
	bad := edge.NewEndpointByProperty([]edge.PropertyMatcher{{Operator: "equals", Value: "x"}}, "")
	if err := bad.Validate(); err == nil {
		t.Error("expected error for property matcher with empty key")
	}
	// valid endpoints
	if err := edge.NewEndpointByName("alice", "User").Validate(); err != nil {
		t.Errorf("unexpected error for valid name endpoint: %v", err)
	}
}

func TestEdgeEqualAcrossMatchStrategies(t *testing.T) {
	byID, _ := edge.NewEdge("a", "b", "K", nil)
	byName, _ := edge.NewEdgeWithEndpoints(edge.NewEndpointByName("a", ""), edge.NewEndpointByID("b"), "K", nil)
	if byID.Equal(byName) {
		t.Error("edges differing only in start match strategy must not be equal")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && len(s) > len(substr) && s[len(s)-1] != substr[0]
}
