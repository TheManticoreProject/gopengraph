package gopengraph_test

import (
	"testing"

	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

func TestNewOpenGraph(t *testing.T) {
	g := gopengraph.NewOpenGraph("test")
	if g == nil {
		t.Error("Expected non-nil OpenGraph")
	}
}

func TestAddNode(t *testing.T) {
	g := gopengraph.NewOpenGraph("test")
	n, err := node.NewNode("test-id", []string{"test-label"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	// Test adding new node
	if !g.AddNode(n) {
		t.Error("Expected AddNode to return true for new node")
	}

	// Test adding duplicate node
	if g.AddNode(n) {
		t.Error("Expected AddNode to return false for duplicate node")
	}
}

func TestFindPaths(t *testing.T) {
	g := gopengraph.NewOpenGraph("test")

	// Create test nodes
	n1, err := node.NewNode("1", []string{"node1"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	n2, err := node.NewNode("2", []string{"node2"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	n3, err := node.NewNode("3", []string{"node3"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	g.AddNode(n1)
	g.AddNode(n2)
	g.AddNode(n3)

	// Create edges
	e1, err := edge.NewEdge("1", "2", "CONNECTS_TO", properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create edge: %v", err)
	}
	e2, err := edge.NewEdge("2", "3", "CONNECTS_TO", properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create edge: %v", err)
	}

	g.AddEdge(e1)
	g.AddEdge(e2)

	// Test path finding
	paths := g.FindPaths("1", "3", 2)
	if len(paths) != 1 {
		t.Errorf("Expected 1 path, got %d", len(paths))
	}
	if len(paths[0]) != 3 {
		t.Errorf("Expected path length 3, got %d", len(paths[0]))
	}

	// Test non-existent path
	paths = g.FindPaths("1", "4", 2)
	if paths != nil {
		t.Error("Expected nil paths for non-existent node")
	}
}

func TestGetConnectedComponents(t *testing.T) {
	g := gopengraph.NewOpenGraph("test")

	// Create two separate components
	n1, err := node.NewNode("1", []string{"node1"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	n2, err := node.NewNode("2", []string{"node2"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	n3, err := node.NewNode("3", []string{"node3"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}
	n4, err := node.NewNode("4", []string{"node4"}, properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create node: %v", err)
	}

	g.AddNode(n1)
	g.AddNode(n2)
	g.AddNode(n3)
	g.AddNode(n4)

	e1, err := edge.NewEdge("1", "2", "CONNECTS_TO", properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create edge: %v", err)
	}
	e2, err := edge.NewEdge("3", "4", "CONNECTS_TO", properties.NewProperties())
	if err != nil {
		t.Fatalf("Failed to create edge: %v", err)
	}

	g.AddEdge(e1)
	g.AddEdge(e2)

	components := g.GetConnectedComponents()
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Check component sizes
	for _, comp := range components {
		if len(comp) != 2 {
			t.Errorf("Expected component size 2, got %d", len(comp))
		}
	}
}
