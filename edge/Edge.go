package edge

import (
	"fmt"

	"github.com/TheManticoreProject/gopengraph/properties"
)

// Edge represents a directed edge in the OpenGraph.
// Follows BloodHound OpenGraph schema requirements with start/end nodes, kind, and properties.
// All edges are directed and one-way as per BloodHound requirements.
//
// Sources:
// - https://bloodhound.specterops.io/opengraph/schema#edges
// - https://bloodhound.specterops.io/opengraph/schema#minimal-working-json
type Edge struct {
	startNodeID string
	endNodeID   string
	kind        string
	properties  *properties.Properties
}

// NewEdge creates a new Edge instance
func NewEdge(startNodeID string, endNodeID string, kind string, p *properties.Properties) (*Edge, error) {
	if startNodeID == "" {
		return nil, fmt.Errorf("start node ID cannot be empty")
	}
	if endNodeID == "" {
		return nil, fmt.Errorf("end node ID cannot be empty")
	}
	if kind == "" {
		return nil, fmt.Errorf("edge kind cannot be empty")
	}

	if p == nil {
		p = properties.NewProperties()
	}

	return &Edge{
		startNodeID: startNodeID,
		endNodeID:   endNodeID,
		kind:        kind,
		properties:  p,
	}, nil
}

// SetProperty sets a property on the edge
func (e *Edge) SetProperty(key string, value interface{}) {
	e.properties.SetProperty(key, value)
}

// GetProperty gets a property from the edge
func (e *Edge) GetProperty(key string, defaultVal ...interface{}) interface{} {
	return e.properties.GetProperty(key, defaultVal...)
}

// GetProperties returns the properties of the edge
func (e *Edge) GetProperties() *properties.Properties {
	return e.properties
}

// RemoveProperty removes a property from the edge
func (e *Edge) RemoveProperty(key string) {
	e.properties.RemoveProperty(key)
}

// ToDict converts edge to map for JSON serialization
func (e *Edge) ToDict() map[string]interface{} {
	edgeDict := map[string]interface{}{
		"kind": e.kind,
		"start": map[string]string{
			"value":    e.startNodeID,
			"match_by": "id",
		},
		"end": map[string]string{
			"value":    e.endNodeID,
			"match_by": "id",
		},
	}

	// Only include properties if they exist and are not empty
	if props := e.properties.ToDict(); len(props) > 0 {
		edgeDict["properties"] = props
	}

	return edgeDict
}

// GetStartNodeID returns the start node ID
func (e *Edge) GetStartNodeID() string {
	return e.startNodeID
}

// GetEndNodeID returns the end node ID
func (e *Edge) GetEndNodeID() string {
	return e.endNodeID
}

// GetKind returns the edge kind/type
func (e *Edge) GetKind() string {
	return e.kind
}

// Equal checks if two edges are equal based on their start, end, and kind
func (e *Edge) Equal(other *Edge) bool {
	if other == nil {
		return false
	}
	return e.startNodeID == other.startNodeID &&
		e.endNodeID == other.endNodeID &&
		e.kind == other.kind
}

// String returns a string representation of the edge
func (e *Edge) String() string {
	return fmt.Sprintf("Edge(start='%s', end='%s', kind='%s', properties=%v)",
		e.startNodeID, e.endNodeID, e.kind, e.properties.ToDict())
}
