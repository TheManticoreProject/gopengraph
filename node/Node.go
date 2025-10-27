package node

import (
	"fmt"

	"github.com/TheManticoreProject/gopengraph/properties"
)

// Node represents a node in the OpenGraph.
// Follows BloodHound OpenGraph schema requirements with unique IDs, kinds, and properties.
//
// Sources:
// - https://bloodhound.specterops.io/opengraph/schema#nodes
// - https://bloodhound.specterops.io/opengraph/schema#minimal-working-json
type Node struct {
	id         string
	kinds      []string
	properties *properties.Properties
}

// NewNode creates a new Node instance
func NewNode(id string, kinds []string, p *properties.Properties) (*Node, error) {
	if id == "" {
		return nil, fmt.Errorf("node ID cannot be empty")
	}

	if kinds == nil {
		kinds = make([]string, 0)
	}

	if p == nil {
		p = properties.NewProperties()
	}

	return &Node{
		id:         id,
		kinds:      kinds,
		properties: p,
	}, nil
}

// AddKind adds a kind/type to the node if it doesn't already exist
func (n *Node) AddKind(kind string) {
	if !n.HasKind(kind) {
		n.kinds = append(n.kinds, kind)
	}
}

// RemoveKind removes a kind/type from the node if it exists
func (n *Node) RemoveKind(kind string) {
	for i, k := range n.kinds {
		if k == kind {
			n.kinds = append(n.kinds[:i], n.kinds[i+1:]...)
			return
		}
	}
}

func (n *Node) GetKinds() []string {
	return n.kinds
}

// HasKind checks if node has a specific kind/type
func (n *Node) HasKind(kind string) bool {
	for _, k := range n.kinds {
		if k == kind {
			return true
		}
	}
	return false
}

func (n *Node) GetID() string {
	return n.id
}

// SetProperty sets a property on the node
func (n *Node) SetProperty(key string, value interface{}) {
	n.properties.SetProperty(key, value)
}

// GetProperty gets a property from the node
func (n *Node) GetProperty(key string, defaultVal ...interface{}) interface{} {
	return n.properties.GetProperty(key, defaultVal...)
}

// GetProperties returns the properties of the node
func (n *Node) GetProperties() *properties.Properties {
	return n.properties
}

// RemoveProperty removes a property from the node
func (n *Node) RemoveProperty(key string) {
	n.properties.RemoveProperty(key)
}

// ToDict converts node to map for JSON serialization
func (n *Node) ToDict() map[string]interface{} {
	return map[string]interface{}{
		"id":         n.id,
		"kinds":      append([]string{}, n.kinds...),
		"properties": n.properties.GetAllProperties(),
	}
}

// Equal checks if two nodes are equal based on their ID
func (n *Node) Equal(other *Node) bool {
	if other == nil {
		return false
	}
	return n.id == other.id
}

// String returns a string representation of the Node
func (n *Node) String() string {
	return fmt.Sprintf("Node(id='%s', kinds=%v, properties=%v)", n.id, n.kinds, n.properties.ToDict())
}
