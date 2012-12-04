Algorithms
==========

In this repository, some algorithms are implemented in go language.

maxflow
-------

Package maxflow implements the max-flow(min-cuts) algorithm which is used in computer vision.

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