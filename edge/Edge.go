package edge

import (
	"fmt"

	"github.com/TheManticoreProject/gopengraph/properties"
)

// match_by strategies for resolving an edge endpoint to a node, as defined by
// the BloodHound OpenGraph edge schema.
//
// Source: https://bloodhound.specterops.io/opengraph/developer/edges
const (
	// MatchByID resolves an endpoint by unique node id. This is the default.
	MatchByID = "id"
	// MatchByName resolves an endpoint by name. Deprecated in BloodHound but
	// still accepted; supported here so payloads using it can round-trip.
	MatchByName = "name"
	// MatchByProperty resolves an endpoint dynamically from one or more
	// property matchers evaluated at ingestion time.
	MatchByProperty = "property"
)

// PropertyMatcher is a single property-based match criterion, used when an
// endpoint's match strategy is MatchByProperty. Values are restricted to
// primitives (string, number, boolean) by the schema.
type PropertyMatcher struct {
	Key      string
	Operator string
	Value    interface{}
}

// Endpoint identifies the start or end of an edge and describes how BloodHound
// should resolve it to a node.
type Endpoint struct {
	matchBy          string
	value            string
	kind             string
	propertyMatchers []PropertyMatcher
}

// NewEndpointByID creates an endpoint resolved by node id.
func NewEndpointByID(value string) Endpoint {
	return Endpoint{matchBy: MatchByID, value: value}
}

// NewEndpointByName creates an endpoint resolved by name. kind is optional and
// may be empty; when set it disambiguates the kind of the target node.
func NewEndpointByName(value string, kind string) Endpoint {
	return Endpoint{matchBy: MatchByName, value: value, kind: kind}
}

// NewEndpointByProperty creates an endpoint resolved by property matchers. kind
// is optional and may be empty.
func NewEndpointByProperty(matchers []PropertyMatcher, kind string) Endpoint {
	return Endpoint{matchBy: MatchByProperty, kind: kind, propertyMatchers: matchers}
}

// GetMatchBy returns the endpoint match strategy, defaulting to MatchByID when
// unset.
func (e Endpoint) GetMatchBy() string {
	if e.matchBy == "" {
		return MatchByID
	}
	return e.matchBy
}

// GetValue returns the endpoint value (the id or name). It is empty for
// property-matched endpoints.
func (e Endpoint) GetValue() string { return e.value }

// GetKind returns the optional kind hint for the endpoint.
func (e Endpoint) GetKind() string { return e.kind }

// GetPropertyMatchers returns the property matchers for a property-matched
// endpoint.
func (e Endpoint) GetPropertyMatchers() []PropertyMatcher { return e.propertyMatchers }

// Validate reports whether the endpoint is internally consistent with its match
// strategy.
func (e Endpoint) Validate() error {
	switch e.GetMatchBy() {
	case MatchByID, MatchByName:
		if e.value == "" {
			return fmt.Errorf("endpoint with match_by %q requires a non-empty value", e.GetMatchBy())
		}
	case MatchByProperty:
		if len(e.propertyMatchers) == 0 {
			return fmt.Errorf("endpoint with match_by %q requires at least one property matcher", MatchByProperty)
		}
		for _, m := range e.propertyMatchers {
			if m.Key == "" {
				return fmt.Errorf("property matcher requires a non-empty key")
			}
		}
	default:
		return fmt.Errorf("unsupported match_by %q", e.matchBy)
	}
	return nil
}

// ToDict converts the endpoint to a map for JSON serialization, matching the
// shape expected by BloodHound for the endpoint's match strategy.
func (e Endpoint) ToDict() map[string]interface{} {
	dict := map[string]interface{}{
		"match_by": e.GetMatchBy(),
	}

	if e.GetMatchBy() == MatchByProperty {
		matchers := make([]map[string]interface{}, 0, len(e.propertyMatchers))
		for _, m := range e.propertyMatchers {
			matchers = append(matchers, map[string]interface{}{
				"key":      m.Key,
				"operator": m.Operator,
				"value":    m.Value,
			})
		}
		dict["property_matchers"] = matchers
	} else {
		dict["value"] = e.value
	}

	if e.kind != "" {
		dict["kind"] = e.kind
	}

	return dict
}

// Equal reports whether two endpoints are equivalent.
func (e Endpoint) Equal(other Endpoint) bool {
	if e.GetMatchBy() != other.GetMatchBy() || e.value != other.value || e.kind != other.kind {
		return false
	}
	if len(e.propertyMatchers) != len(other.propertyMatchers) {
		return false
	}
	for i := range e.propertyMatchers {
		if e.propertyMatchers[i] != other.propertyMatchers[i] {
			return false
		}
	}
	return true
}

// Edge represents a directed edge in the OpenGraph.
// Follows BloodHound OpenGraph schema requirements with start/end endpoints,
// kind, and properties. All edges are directed and one-way as per BloodHound
// requirements.
//
// Sources:
// - https://bloodhound.specterops.io/opengraph/developer/edges
// - https://bloodhound.specterops.io/opengraph/developer/graph-data
type Edge struct {
	start      Endpoint
	end        Endpoint
	kind       string
	properties *properties.Properties
}

// NewEdge creates a new Edge instance whose endpoints are resolved by node id.
// This is the common case; use NewEdgeWithEndpoints for name or property
// matching.
func NewEdge(startNodeID string, endNodeID string, kind string, p *properties.Properties) (*Edge, error) {
	if startNodeID == "" {
		return nil, fmt.Errorf("start node ID cannot be empty")
	}
	if endNodeID == "" {
		return nil, fmt.Errorf("end node ID cannot be empty")
	}

	return NewEdgeWithEndpoints(NewEndpointByID(startNodeID), NewEndpointByID(endNodeID), kind, p)
}

// NewEdgeWithEndpoints creates a new Edge instance from explicit endpoints,
// allowing any match strategy for either end.
func NewEdgeWithEndpoints(start Endpoint, end Endpoint, kind string, p *properties.Properties) (*Edge, error) {
	if kind == "" {
		return nil, fmt.Errorf("edge kind cannot be empty")
	}
	if err := start.Validate(); err != nil {
		return nil, fmt.Errorf("invalid start endpoint: %w", err)
	}
	if err := end.Validate(); err != nil {
		return nil, fmt.Errorf("invalid end endpoint: %w", err)
	}

	if p == nil {
		p = properties.NewProperties()
	}

	return &Edge{
		start:      start,
		end:        end,
		kind:       kind,
		properties: p,
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
		"kind":  e.kind,
		"start": e.start.ToDict(),
		"end":   e.end.ToDict(),
	}

	// Only include properties if they exist and are not empty
	if props := e.properties.ToDict(); len(props) > 0 {
		edgeDict["properties"] = props
	}

	return edgeDict
}

// GetStart returns the start endpoint.
func (e *Edge) GetStart() Endpoint { return e.start }

// GetEnd returns the end endpoint.
func (e *Edge) GetEnd() Endpoint { return e.end }

// GetStartNodeID returns the start endpoint value. For id-matched endpoints this
// is the node id; it is empty for property-matched endpoints.
func (e *Edge) GetStartNodeID() string {
	return e.start.value
}

// GetEndNodeID returns the end endpoint value. For id-matched endpoints this is
// the node id; it is empty for property-matched endpoints.
func (e *Edge) GetEndNodeID() string {
	return e.end.value
}

// GetKind returns the edge kind/type
func (e *Edge) GetKind() string {
	return e.kind
}

// Equal checks if two edges are equal based on their endpoints and kind
func (e *Edge) Equal(other *Edge) bool {
	if other == nil {
		return false
	}
	return e.kind == other.kind &&
		e.start.Equal(other.start) &&
		e.end.Equal(other.end)
}

// String returns a string representation of the edge
func (e *Edge) String() string {
	return fmt.Sprintf("Edge(start='%s', end='%s', kind='%s', properties=%v)",
		e.start.value, e.end.value, e.kind, e.properties.ToDict())
}
