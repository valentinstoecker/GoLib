package graph

import "errors"

type graph struct {
	dPCreate func() DataProvider
	edges    map[EdgeID]edge
	nodes    map[NodeID]node
}

func (g *graph) AddEdge(src, tgt NodeID, data interface{}) (EdgeID, error) {
	srcNode, ok := g.nodes[src]
	if !ok {
		return 0, errors.New("Source doesn't exist")
	}
	tgtNode, ok := g.nodes[tgt]
	if !ok {
		return 0, errors.New("Target doesn't exist")
	}
	currEdgeID++
	srcNode.es = append(srcNode.es, currEdgeID)
	tgtNode.es = append(tgtNode.es, currEdgeID)
	dp := g.dPCreate()
	dp.PutData(data)
	e := edge{
		DataProvider: dp,
		id:           currEdgeID,
		src:          src,
		tgt:          tgt,
		g:            g,
	}
	g.edges[currEdgeID] = e
	g.nodes[src] = srcNode
	g.nodes[tgt] = tgtNode
	return currEdgeID, nil
}

func (g *graph) AddNode(data interface{}) (NodeID, error) {
	currNodeID++
	dp := g.dPCreate()
	dp.PutData(data)
	n := node{
		DataProvider: dp,
		id:           currNodeID,
		g:            g,
	}
	g.nodes[currNodeID] = n
	return currNodeID, nil
}

func (g *graph) Edge(id EdgeID) (Edge, error) {
	if e, ok := g.edges[id]; ok {
		return e, nil
	}
	return edge{}, errors.New("Didn't find Edge")
}

func (g *graph) Node(id NodeID) (Node, error) {
	if n, ok := g.nodes[id]; ok {
		return n, nil
	}
	return node{}, errors.New("Didn't find Node")
}

func test() Graph {
	return &graph{}
}
