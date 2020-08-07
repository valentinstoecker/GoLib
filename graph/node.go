package graph

var currNodeID NodeID

type node struct {
	DataProvider
	id NodeID
	es []EdgeID
	g  *graph
}

func (n node) ID() NodeID {
	return n.id
}

func (n node) Edges() []Edge {
	es := make([]Edge, 0, len(n.es))
	for _, eID := range n.es {
		if e, err := (*n.g).Edge(eID); err == nil {
			es = append(es, e)
		}
	}
	return es
}
