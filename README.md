Algorithms
==========

In this repository, some algorithms are implemented in go language.

GoDoc link: [ed](http://godoc.org/github.com/daviddengcn/go-algs/ed) [maxflow](http://godoc.org/github.com/daviddengcn/go-algs/maxflow)

maxflow
-------

This package implements the max-flow(min-cuts, graph-cuts) algorithm that is used to solve the labeling problem in computer vision or graphics area.

Usage:
```go
    g := maxflow.NewGraph()
    
    nodes := make([]*maxflow.Node, 2)
    
    for i := range(nodes) {
        nodes[i] = g.AddNode()
    } // for i
    
    g.SetTweights(nodes[0], 1, 5)
    g.SetTweights(nodes[1], 2, 6)
    g.AddEdge(nodes[0], nodes[1], 3, 4)
    
    g.Run();

    flow := g.Flow()

    isSource0 := g.IsSource(nodes[0])
```

ed
--------

This package implements the edit-distance algorithm that is used to compute the similarity between two strings, or more generally defined lists.

1. For computing the standard edit-distance of two strings, call ed.String or ed.StringFull function.

1. For generally defined lists, implement the ed.Interface, and use ed.EditDistance or ed.EditDistanceFull function.


LICENSE
-------
BSD license