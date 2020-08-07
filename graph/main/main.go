package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/valentinstoecker/GoLib/graph"
)

type person struct {
	Name string
	Age  int
}

type parent struct {
}

type worksFor struct {
	start time.Time
}

func (w worksFor) String() string {
	return fmt.Sprintf("Since%s", w.start.Format("_2 Jan 2006"))
}

func main() {
	g := graph.StdGraph()
	f, _ := g.AddNode(person{
		Name: "Frank Stöcker",
		Age:  60,
	})
	c, _ := g.AddNode(person{
		Name: "Christa Stöcker-Lenz",
		Age:  57,
	})
	v, _ := g.AddNode(person{
		Name: "Valentin Stöcker",
		Age:  19,
	})
	m, _ := g.AddNode(person{
		Name: "Moritz Staffel",
		Age:  27,
	})
	g.AddEdge(f, v, parent{})
	g.AddEdge(c, v, parent{})
	g.AddEdge(v, m, worksFor{
		start: time.Date(2017, time.October, 4, 0, 0, 0, 0, time.UTC),
	})

	nv, _ := g.Node(v)
	fmt.Println(graph.EdgesOfType(nv, reflect.TypeOf(parent{})))
	fmt.Println(graph.EdgesOfType(nv, reflect.TypeOf(worksFor{}))[0].Data())
}
