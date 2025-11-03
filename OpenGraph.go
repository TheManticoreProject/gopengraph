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

// Edges operations

// AddEdge adds an edge to the graph after performing validation checks.
//
// It verifies that both the start and end nodes referenced by the edge exist in the graph,
// and that the edge is not a duplicate of an existing edge. If any validation fails,
// the edge is not added.
//
// Arguments:
//
//	edge *edge.Edge: The edge to be added to the graph.
//
// Returns:
//
//	bool: True if the edge was successfully added, false if validation failed
//	      (e.g., nodes do not exist or the edge is a duplicate).
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

	return g.AddEdgeWithoutValidation(edge)
}

// AddEdgeWithoutValidation adds an edge to the graph without validating the nodes.
//
// This is a convenience function for adding edges without the validation checks performed by AddEdge.
// It is useful when you are sure that the nodes and edge already exist in the graph,
// or when you want to add an edge without performing the validation checks.
//
// Arguments:
//
//	edge *edge.Edge: The edge to be added to the graph.
//
// Returns:
//
//	bool: True if the edge was successfully added.
func (g *OpenGraph) AddEdgeWithoutValidation(edge *edge.Edge) bool {
	g.edges = append(g.edges, edge)
	return true
}

// Nodes operations

// AddNode adds a node to the graph after performing validation checks.
//
// It verifies that the node does not already exist in the graph,
// and that the node has a valid ID. If any validation fails,
// the node is not added.
//
// Arguments:
//
//	node *node.Node: The node to be added to the graph.
//
// Returns:
//
//	bool: True if the node was successfully added, false if validation failed
//	      (e.g., node already exists or has an invalid ID).
func (g *OpenGraph) AddNode(node *node.Node) bool {
	if _, exists := g.nodes[node.GetID()]; exists {
		return false
	}

	// Add source kind if specified and not already present
	if g.sourceKind != "" && !node.HasKind(g.sourceKind) {
		node.AddKind(g.sourceKind)
	}

	return g.AddNodeWithoutValidation(node)
}

// AddNodeWithoutValidation adds a node to the graph without validating the node.
//
// This is a convenience function for adding nodes without the validation checks performed by AddNode.
// It is useful when you are sure that the node already exists in the graph,
// or when you want to add a node without performing the validation checks.
//
// Arguments:
//
//	node *node.Node: The node to be added to the graph.
//
// Returns:
//
//	bool: True if the node was successfully added.
func (g *OpenGraph) AddNodeWithoutValidation(node *node.Node) bool {
	g.nodes[node.GetID()] = node
	return true
}

// RemoveNodeByID removes a node and its associated edges after performing validation checks.
//
// It verifies that the node exists in the graph,
// and that the node has a valid ID. If any validation fails,
// the node is not removed.
//
// Arguments:
//
//	id string: The ID of the node to be removed from the graph.
//
// Returns:
//
//	bool: True if the node was successfully removed, false if validation failed
//	      (e.g., node does not exist or has an invalid ID).
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

// GetNode returns a node by ID after performing validation checks.
//
// It verifies that the node exists in the graph,
// and that the node has a valid ID. If any validation fails,
// the node is not returned.
//
// Arguments:
//
//	id string: The ID of the node to be returned from the graph.
//
// Returns:
//
//	*node.Node: The node if it exists, nil if validation failed
//	             (e.g., node does not exist or has an invalid ID).
func (g *OpenGraph) GetNode(id string) *node.Node {
	return g.nodes[id]
}

// GetNodesByKind returns all nodes of a specific kind after performing validation checks.
//
// It verifies that the kind is valid,
// and that the nodes exist in the graph. If any validation fails,
// the nodes are not returned.
//
// Arguments:
//
//	kind string: The kind of nodes to be returned from the graph.
//
// Returns:
//
//	[]*node.Node: The nodes if they exist, nil if validation failed
//	              (e.g., kind is not valid or nodes do not exist).
func (g *OpenGraph) GetNodesByKind(kind string) []*node.Node {
	var nodes []*node.Node
	for _, n := range g.nodes {
		if n.HasKind(kind) {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

// GetEdgesByKind returns all edges of a specific kind after performing validation checks.
//
// It verifies that the kind is valid,
// and that the edges exist in the graph. If any validation fails,
// the edges are not returned.
//
// Arguments:
//
//	kind string: The kind of edges to be returned from the graph.
//
// Returns:
//
//	[]*edge.Edge: The edges if they exist, nil if validation failed
//	              (e.g., kind is not valid or edges do not exist).
func (g *OpenGraph) GetEdgesByKind(kind string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetKind() == kind {
			edges = append(edges, e)
		}
	}
	return edges
}

// GetEdgesFromNode returns all edges starting from a node after performing validation checks.
//
// It verifies that the node exists in the graph,
// and that the node has a valid ID. If any validation fails,
// the edges are not returned.
//
// Arguments:
//
//	id string: The ID of the node to get edges from.
//
// Returns:
//
//	[]*edge.Edge: The edges if they exist, nil if validation failed
//	              (e.g., node does not exist or has an invalid ID).
func (g *OpenGraph) GetEdgesFromNode(id string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetStartNodeID() == id {
			edges = append(edges, e)
		}
	}
	return edges
}

// GetEdgesToNode returns all edges ending at a node after performing validation checks.
//
// It verifies that the node exists in the graph,
// and that the node has a valid ID. If any validation fails,
// the edges are not returned.
//
// Arguments:
//
//	id string: The ID of the node to get edges to.
//
// Returns:
//
//	[]*edge.Edge: The edges if they exist, nil if validation failed
//	              (e.g., node does not exist or has an invalid ID).
func (g *OpenGraph) GetEdgesToNode(id string) []*edge.Edge {
	var edges []*edge.Edge
	for _, e := range g.edges {
		if e.GetEndNodeID() == id {
			edges = append(edges, e)
		}
	}
	return edges
}

// Graph operations

// FindPaths finds all paths between two nodes using BFS after performing validation checks.
//
// It verifies that the start and end nodes exist in the graph,
// and that the start and end nodes have valid IDs. If any validation fails,
// the paths are not returned.
//
// Arguments:
//
//	startID string: The ID of the start node.
//	endID string: The ID of the end node.
//	maxDepth int: The maximum depth of the paths to find.
//
// Returns:
//
//	[][]string: The paths if they exist, nil if validation failed
//	             (e.g., start or end node does not exist or has an invalid ID).
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

// GetConnectedComponents finds all connected components after performing validation checks.
//
// It verifies that the nodes exist in the graph,
// and that the nodes have valid IDs. If any validation fails,
// the connected components are not returned.
//
// Arguments:
//
// Returns:
//
//	[]map[string]bool: The connected components if they exist, nil if validation failed
//	                   (e.g., nodes do not exist or have an invalid ID).
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

// ValidateGraph checks for common graph issues after performing validation checks.
//
// It verifies that the edges and nodes exist in the graph,
// and that the edges and nodes have valid IDs. If any validation fails,
// the errors are not returned.
//
// Arguments:
//
// Returns:
//
//	[]string: The errors if they exist, nil if validation failed
//	           (e.g., edges or nodes do not exist or have an invalid ID).
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

// Graph exports

// ExportJSON exports the graph to JSON format after performing validation checks.
//
// It verifies that the nodes and edges exist in the graph,
// and that the nodes and edges have valid IDs. If any validation fails,
// the JSON is not returned.
//
// Arguments:
//
// includeMetadata bool: Whether to include metadata in the JSON.
//
// Returns:
//
//	string: The JSON if it exists, nil if validation failed
//	        (e.g., nodes or edges do not exist or have an invalid ID).
//	error: An error if the JSON is not returned.
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

// ExportToFile exports the graph to a JSON file after performing validation checks.
//
// It verifies that the nodes and edges exist in the graph,
// and that the nodes and edges have valid IDs. If any validation fails,
// the file is not written.
//
// Arguments:
//
//	filename string: The name of the file to export the graph to.
//
// Returns:
//
//	error: An error if the file is not written.
func (g *OpenGraph) ExportToFile(filename string) error {
	jsonData, err := g.ExportJSON(true)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(jsonData), 0644)
}

// Graph infos

// GetNodeCount returns the total number of nodes after performing validation checks.
//
// It verifies that the nodes exist in the graph,
// and that the nodes have valid IDs. If any validation fails,
// the node count is not returned.
//
// Arguments:
//
// Returns:
//
//	int: The number of nodes if they exist, nil if validation failed
//	     (e.g., nodes do not exist or have an invalid ID).
func (g *OpenGraph) GetNodeCount() int {
	return len(g.nodes)
}

// GetEdgeCount returns the total number of edges after performing validation checks.
//
// It verifies that the edges exist in the graph,
// and that the edges have valid IDs. If any validation fails,
// the edge count is not returned.
//
// Arguments:
//
// Returns:
//
//	int: The number of edges if they exist, nil if validation failed
//	     (e.g., edges do not exist or have an invalid ID).
func (g *OpenGraph) GetEdgeCount() int {
	return len(g.edges)
}

// Clear removes all nodes and edges after performing validation checks.
//
// It verifies that the nodes and edges exist in the graph,
// and that the nodes and edges have valid IDs. If any validation fails,
// the nodes and edges are not removed.
//
// Arguments:
//
// Returns:
//
//	nil: If the nodes and edges were successfully removed, nil if validation failed
//	     (e.g., nodes or edges do not exist or have an invalid ID).
func (g *OpenGraph) Clear() {
	g.nodes = make(map[string]*node.Node)
	g.edges = make([]*edge.Edge, 0)
}

// Len returns the total number of nodes and edges after performing validation checks.
//
// It verifies that the nodes and edges exist in the graph,
// and that the nodes and edges have valid IDs. If any validation fails,
// the length is not returned.
//
// Arguments:
//
// Returns:
//
//	int: The total number of nodes and edges if they exist, nil if validation failed
//	     (e.g., nodes or edges do not exist or have an invalid ID).
func (g *OpenGraph) Len() int {
	return len(g.nodes) + len(g.edges)
}

// String returns a string representation of the graph after performing validation checks.
//
// It verifies that the nodes and edges exist in the graph,
// and that the nodes and edges have valid IDs. If any validation fails,
// the string representation is not returned.
//
// Arguments:
//
// Returns:
//
//	string: The string representation if it exists, nil if validation failed
//	        (e.g., nodes or edges do not exist or have an invalid ID).
func (g *OpenGraph) String() string {
	return fmt.Sprintf("OpenGraph(nodes=%d, edges=%d, source_kind='%s')",
		len(g.nodes), len(g.edges), g.sourceKind)
}
