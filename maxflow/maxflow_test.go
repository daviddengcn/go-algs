package maxflow

import (
	"testing"
)

func assertResults(t *testing.T, expFlow CapType, g *Graph, nodes []*Node, isSources []bool) {
	if g.Flow() != expFlow {
		t.Errorf("Maxflow should be %d but get %d", expFlow, g.Flow())
	} // if

	for i := range nodes {
		if g.IsSource(nodes[i]) != isSources[i] {
			if isSources[i] {
				t.Errorf("node %d should be SOURCE!", i)
			} else {
				t.Errorf("node %d should be SINK!", i)
			} // else
		} // if
	} // for i
}

func TestMaxflow1(t *testing.T) {
	g := NewGraph()

	nodes := make([]*Node, 2)

	for i := range nodes {
		nodes[i] = g.AddNode()
	} // for i

	g.SetTweights(nodes[0], 1, 5)
	g.SetTweights(nodes[1], 2, 6)
	g.AddEdge(nodes[0], nodes[1], 3, 4)

	g.Run()
	assertResults(t, 3, g, nodes, []bool{false, false})
}

func TestMaxflow2(t *testing.T) {
	g := NewGraph()

	nodes := make([]*Node, 4)

	for i := range nodes {
		nodes[i] = g.AddNode()
	} // for i

	g.SetTweights(nodes[0], 3, 0)
	g.SetTweights(nodes[1], 3, 0)
	g.SetTweights(nodes[2], 0, 2)
	g.SetTweights(nodes[3], 0, 3)

	g.AddEdge(nodes[0], nodes[1], 2, 0)
	g.AddEdge(nodes[0], nodes[2], 3, 0)
	g.AddEdge(nodes[1], nodes[3], 2, 0)
	g.AddEdge(nodes[2], nodes[3], 4, 0)

	g.Run()
	assertResults(t, 5, g, nodes, []bool{true, true, false, false})
}
