package main

import (
	"github.com/TheManticoreProject/gopengraph"
	"github.com/TheManticoreProject/gopengraph/edge"
	"github.com/TheManticoreProject/gopengraph/node"
	"github.com/TheManticoreProject/gopengraph/properties"
)

func main() {
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

	// Export to file
	graph.ExportToFile("minimal_working_json.json")
}
