package graph

import "fmt"

var currEdgeID EdgeID

type edge struct {
	DataProvider
	id  EdgeID
	src NodeID
	tgt NodeID
	g   *graph
}

func (e edge) ID() EdgeID {
	return e.id
}

func (e edge) Source() NodeID {
	return e.src
}

func (e edge) Target() NodeID {
	return e.tgt
}

func (e edge) String() string {
	return fmt.Sprintf("Edge(%d) %d -> %d Data: %v", e.id, e.src, e.tgt, e.Data())
}
