package gopengraph

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
)

// OpenGraph struct for managing a graph structure compatible with BloodHound OpenGraph.
//
// Follows BloodHound OpenGraph schema requirements and best practices.
//
// Sources:
// - https://bloodhound.specterops.io/opengraph/schema#opengraph
// - https://bloodhound.specterops.io/opengraph/schema#minimal-working-json
// - https://bloodhound.specterops.io/opengraph/best-practices
type OpenGraph struct {
	nodes      map[string]*node.Node
	edges      []*edge.Edge
	sourceKind string
}

// NewOpenGraph creates a new OpenGraph instance
func NewOpenGraph(sourceKind string) *OpenGraph {
	return &OpenGraph{
		nodes:      make(map[string]*node.Node),
		edges:      make([]*edge.Edge, 0),
		sourceKind: sourceKind,
	}
}

// AddNode adds a node to the graph
func (g *OpenGraph) AddNode(node *node.Node) bool {
	if _, exists := g.nodes[node.GetID()]; exists {
		return false
	}

	// Add source kind if specified and not already present
	if g.sourceKind != "" && !node.HasKind(g.sourceKind) {
		node.AddKind(g.sourceKind)
	}

	g.nodes[node.GetID()] = node
	return true
}

// AddEdge adds an edge to the graph
func (g *OpenGraph) AddEdge(edge *edge.Edge) bool {
	// Verify both nodes exist
	if _, exists := g.nodes[edge.GetStartNodeID()]; !exists {
		return false
	}
	if _, exists := g.nodes[edge.GetEndNodeID()]; !exists {
		return false
	}

	// Check for duplicate edge
	for _, e := range g.edges {
		if e.Equal(edge) {
			return false
		}
	}

	g.edges = append(g.edges, edge)
	return true
}

// RemoveNodeByID removes a node and its associated edges
func (g *OpenGraph) RemoveNodeByID(id string) bool {
	if _, exists := g.nodes[id]; !exists {
		return false
	}

	delete(g.nodes, id)

	// Remove associated edges
	newEdges := make([]*edge.Edge, 0)
	for _, e := range g.edges {
		if e.GetStartNodeID() != id && e.GetEndNodeID() != id {
			newEdges = append(newEdges, e)
		}
	}
	g.edges = newEdges

	return true
}

// GetNode returns a node by ID
func (g *OpenGraph) GetNode(id string) *node.Node {
	return g.nodes[id]
}

// GetNodesByKind returns all nodes of a specific kind
func (g *OpenGraph) GetNodesByKind(kind string) []*node.Node {
	var nodes []*node.Node
	for _, n := range g.nodes {
		if n.HasKind(kind) {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

// GetEdgesByKind returns all edges of a specific kind
func (g *OpenGraph) GetEdgesByKind(kind string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetKind() == kind {
			edges = append(edges, e)
		}
	}
	return edges
}

// GetEdgesFromNode returns all edges starting from a node
func (g *OpenGraph) GetEdgesFromNode(id string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetStartNodeID() == id {
			edges = append(edges, e)
		}
	}
	return edges
}

// GetEdgesToNode returns all edges ending at a node
func (g *OpenGraph) GetEdgesToNode(id string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetEndNodeID() == id {
			edges = append(edges, e)
		}
	}
	return edges
}

// FindPaths finds all paths between two nodes using BFS
func (g *OpenGraph) FindPaths(startID, endID string, maxDepth int) [][]string {
	if _, exists := g.nodes[startID]; !exists {
		return nil
	}
	if _, exists := g.nodes[endID]; !exists {
		return nil
	}

	if startID == endID {
		return [][]string{{startID}}
	}

	var paths [][]string
	visited := make(map[string]bool)
	queue := []struct {
		id   string
		path []string
	}{{startID, []string{startID}}}
	visited[startID] = true

	for len(queue) > 0 && len(queue[0].path) <= maxDepth {
		current := queue[0]
		queue = queue[1:]

		for _, edge := range g.GetEdgesFromNode(current.id) {
			nextID := edge.GetEndNodeID()
			if !visited[nextID] {
				newPath := append([]string{}, current.path...)
				newPath = append(newPath, nextID)

				if nextID == endID {
					paths = append(paths, newPath)
				} else {
					visited[nextID] = true
					queue = append(queue, struct {
						id   string
						path []string
					}{nextID, newPath})
				}
			}
		}
	}

	return paths
}

// GetConnectedComponents finds all connected components
func (g *OpenGraph) GetConnectedComponents() []map[string]bool {
	visited := make(map[string]bool)
	var components []map[string]bool

	for nodeID := range g.nodes {
		if !visited[nodeID] {
			component := make(map[string]bool)
			stack := []string{nodeID}

			for len(stack) > 0 {
				current := stack[len(stack)-1]
				stack = stack[:len(stack)-1]

				if !visited[current] {
					visited[current] = true
					component[current] = true

					// Add adjacent nodes
					for _, edge := range g.GetEdgesFromNode(current) {
						if !visited[edge.GetEndNodeID()] {
							stack = append(stack, edge.GetEndNodeID())
						}
					}
					for _, edge := range g.GetEdgesToNode(current) {
						if !visited[edge.GetStartNodeID()] {
							stack = append(stack, edge.GetStartNodeID())
						}
					}
				}
			}
			components = append(components, component)
		}
	}

	return components
}

// ValidateGraph checks for common graph issues
func (g *OpenGraph) ValidateGraph() []string {
	var errors []string

	// Check for orphaned edges
	for _, edge := range g.edges {
		if _, exists := g.nodes[edge.GetStartNodeID()]; !exists {
			errors = append(errors, fmt.Sprintf("Edge %s references non-existent start node: %s",
				edge.GetKind(), edge.GetStartNodeID()))
		}
		if _, exists := g.nodes[edge.GetEndNodeID()]; !exists {
			errors = append(errors, fmt.Sprintf("Edge %s references non-existent end node: %s",
				edge.GetKind(), edge.GetEndNodeID()))
		}
	}

	// Check for isolated nodes
	var isolatedNodes []string
	for id := range g.nodes {
		if len(g.GetEdgesFromNode(id)) == 0 && len(g.GetEdgesToNode(id)) == 0 {
			isolatedNodes = append(isolatedNodes, id)
		}
	}

	if len(isolatedNodes) > 0 {
		errors = append(errors, fmt.Sprintf("Found %d isolated nodes: %v",
			len(isolatedNodes), isolatedNodes))
	}

	return errors
}

// ExportJSON exports the graph to JSON format
func (g *OpenGraph) ExportJSON(includeMetadata bool) (string, error) {
	graphData := make(map[string]interface{})
	graphContent := make(map[string]interface{})

	// Convert nodes to dict format
	// Initialize nodesData as an empty slice (not nil) so it marshals to [] instead of null if no nodes exist.
	nodesData := make([]map[string]interface{}, 0, len(g.nodes))
	for _, n := range g.nodes {
		nodesData = append(nodesData, n.ToDict())
	}
	graphContent["nodes"] = nodesData

	// Convert edges to dict format
	// Initialize edgesData as an empty slice (not nil) so it marshals to [] instead of null if no edges exist.
	// The original `var edgesData []map[string]interface{}` declares a nil slice.
	// If `g.edges` is empty, the loop is skipped, and `edgesData` remains nil,
	// which `json.Marshal` converts to `null`.
	// By using `make([]map[string]interface{}, 0)`, it's explicitly an empty slice,
	// which `json.Marshal` converts to `[]`.
	edgesData := make([]map[string]interface{}, 0, len(g.edges))
	for _, e := range g.edges {
		edgesData = append(edgesData, e.ToDict())
	}
	graphContent["edges"] = edgesData

	graphData["graph"] = graphContent

	if includeMetadata && g.sourceKind != "" {
		graphData["metadata"] = map[string]interface{}{
			"source_kind": g.sourceKind,
		}
	}

	jsonData, err := json.MarshalIndent(graphData, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// ExportToFile exports the graph to a JSON file
func (g *OpenGraph) ExportToFile(filename string) error {
	jsonData, err := g.ExportJSON(true)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(jsonData), 0644)
}

// GetNodeCount returns the total number of nodes
func (g *OpenGraph) GetNodeCount() int {
	return len(g.nodes)
}

// GetEdgeCount returns the total number of edges
func (g *OpenGraph) GetEdgeCount() int {
	return len(g.edges)
}

// Clear removes all nodes and edges
func (g *OpenGraph) Clear() {
	g.nodes = make(map[string]*node.Node)
	g.edges = make([]*edge.Edge, 0)
}

// Len returns the total number of nodes and edges
func (g *OpenGraph) Len() int {
	return len(g.nodes) + len(g.edges)
}

func (g *OpenGraph) String() string {
	return fmt.Sprintf("OpenGraph(nodes=%d, edges=%d, source_kind='%s')",
		len(g.nodes), len(g.edges), g.sourceKind)
}
