package graph

import (
	"reflect"
)

// NodeID -> ID of a Node
type NodeID uint64

// EdgeID -> ID of a Node
type EdgeID uint64

// Edge -> A Node of the Graph
type Edge interface {
	ID() EdgeID
	Source() NodeID
	Target() NodeID
	DataProvider
}

// Node -> A Node of the Graph
type Node interface {
	ID() NodeID
	Edges() []Edge
	DataProvider
}

// Graph -> A Graph data structure
type Graph interface {
	Node(NodeID) (Node, error)
	Edge(EdgeID) (Edge, error)
	AddNode(interface{}) (NodeID, error)
	AddEdge(src, tgt NodeID, data interface{}) (EdgeID, error)
}

// DataProvider -> A thing that can provide Data
type DataProvider interface {
	Type() reflect.Type
	Data() interface{}
	DataAs(interface{}) bool
	PutData(interface{})
}

// Incoming -> Gets incoming Edges
func Incoming(n Node) []Edge {
	all := n.Edges()
	inc := make([]Edge, 0, len(all))
	for _, e := range all {
		if e.Target() == n.ID() {
			inc = append(inc, e)
		}
	}
	return inc
}

// Outgoing -> Gets outgoing Edges
func Outgoing(n Node) []Edge {
	all := n.Edges()
	inc := make([]Edge, 0, len(all))
	for _, e := range all {
		if e.Source() == n.ID() {
			inc = append(inc, e)
		}
	}
	return inc
}

// EdgesOfType -> Get Edges with specific types
func EdgesOfType(n Node, t reflect.Type) []Edge {
	all := n.Edges()
	typ := make([]Edge, 0, len(all))
	for _, e := range all {
		if e.Type() == t {
			typ = append(typ, e)
		}
	}
	return typ
}

type dataStore struct {
	typ  reflect.Type
	data interface{}
}

func (ds *dataStore) Type() reflect.Type {
	return ds.typ
}

func (ds *dataStore) Data() interface{} {
	return ds.data
}

func (ds *dataStore) DataAs(d interface{}) bool {
	if reflect.TypeOf(d) == reflect.PtrTo(ds.typ) {
		dv := reflect.ValueOf(d)
		dv.Elem().Set(reflect.ValueOf(ds.data))
		return true
	}
	return false
}

func (ds *dataStore) PutData(d interface{}) {
	ds.data = d
	ds.typ = reflect.TypeOf(d)
}

func StdGraph() Graph {
	return &graph{
		dPCreate: func() DataProvider {
			return &dataStore{}
		},
		nodes: make(map[NodeID]node),
		edges: make(map[EdgeID]edge),
	}
}
