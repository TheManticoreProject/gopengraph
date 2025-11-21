package gopengraph_test

import (
	"testing"

	"encoding/json"

	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

func TestExportJSON(t *testing.T) {
	t.Run("empty graph exports empty nodes and edges arrays", func(t *testing.T) {
		g := gopengraph.NewOpenGraph("test-source")
		jsonData, err := g.ExportJSON(true)
		if err != nil {
			t.Fatalf("ExportJSON failed: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		graphContent, ok := result["graph"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected 'graph' key in JSON output")
		}

		// Check for empty nodes array
		nodes, ok := graphContent["nodes"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'nodes' key to be an array")
		}
		if len(nodes) != 0 {
			t.Errorf("Expected 'nodes' array to be empty, got %d elements", len(nodes))
		}

		// Check for empty edges array
		edges, ok := graphContent["edges"].([]interface{})
		if !ok {
			t.Fatalf("Expected 'edges' key to be an array")
		}
		if len(edges) != 0 {
			t.Errorf("Expected 'edges' array to be empty, got %d elements", len(edges))
		}
	})

	t.Run("metadata not present when includeMetadata is false", func(t *testing.T) {
		g := gopengraph.NewOpenGraph("test-source") // Ensure sourceKind is set for potential metadata
		jsonData, err := g.ExportJSON(false)
		if err != nil {
			t.Fatalf("ExportJSON failed: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if _, exists := result["metadata"]; exists {
			t.Error("Expected 'metadata' key not to be present when includeMetadata is false")
		}
	})

	t.Run("metadata present when includeMetadata is true and sourceKind is set", func(t *testing.T) {
		g := gopengraph.NewOpenGraph("test-source")
		jsonData, err := g.ExportJSON(true)
		if err != nil {
			t.Fatalf("ExportJSON failed: %v", err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		metadata, ok := result["metadata"].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected 'metadata' key to be present and a map when includeMetadata is true")
		}

		sourceKind, ok := metadata["source_kind"].(string)
		if !ok || sourceKind != "test-source" {
			t.Errorf("Expected 'source_kind' in metadata to be 'test-source', got %v", metadata["source_kind"])
		}
	})
}

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

func TestJSONioInvolution(t *testing.T) {
	// Create an OpenGraph instance
	graph := gopengraph.NewOpenGraph("Base")

	// Create nodes
	bobProps := properties.NewProperties()
	bobProps.SetProperty("displayname", "bob")
	bobProps.SetProperty("property", "a")
	bobProps.SetProperty("objectid", "123")
	bobProps.SetProperty("name", "BOB")

	bobNode, _ := node.NewNode("123", []string{"Person", "Base"}, bobProps)

	aliceProps := properties.NewProperties()
	aliceProps.SetProperty("displayname", "alice")
	aliceProps.SetProperty("property", "b")
	aliceProps.SetProperty("objectid", "234")
	aliceProps.SetProperty("name", "ALICE")

	aliceNode, _ := node.NewNode("234", []string{"Person", "Base"}, aliceProps)

	// Add nodes to graph
	graph.AddNode(bobNode)
	graph.AddNode(aliceNode)

	// Create edge: Bob knows Alice
	knowsEdge, _ := edge.NewEdge(
		bobNode.GetID(),   // Bob is the start
		aliceNode.GetID(), // Alice is the end
		"Knows",
		nil,
	)

	// Add edge to graph
	graph.AddEdge(knowsEdge)

	// ============================

	// Export to JSON
	jsonData, err := graph.ExportJSON(true)
	if err != nil {
		t.Fatalf("ExportJSON failed: %v", err)
	}

	importedGraph := gopengraph.NewOpenGraph("")
	// Import from JSON
	err = importedGraph.FromJSON(jsonData)
	if err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	// Check that the graph is the same
	if importedGraph.GetNodeCount() != graph.GetNodeCount() {
		t.Errorf("Expected %d nodes, got %d", graph.GetNodeCount(), importedGraph.GetNodeCount())
	}
	if importedGraph.GetEdgeCount() != graph.GetEdgeCount() {
		t.Errorf("Expected %d edges, got %d", graph.GetEdgeCount(), importedGraph.GetEdgeCount())
	}
	if importedGraph.GetSourceKind() != graph.GetSourceKind() {
		t.Errorf("Expected source kind '%s', got %s", graph.GetSourceKind(), importedGraph.GetSourceKind())
	}
	if !importedGraph.Equal(graph) {
		t.Errorf("Expected graphs to be equal")
	}
}
