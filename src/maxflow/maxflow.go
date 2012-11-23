package maxflow

import (
//    "fmt"
)

type CapType int

type NodeBlock struct {
    nodes []*Node
    current int
}

func (blk *NodeBlock) New() *Node {
    nd := &Node{}
    blk.nodes = append(blk.nodes, nd)
    return nd
}

func (blk *NodeBlock) scanFirst() *Node {
    blk.current = 0
    if len(blk.nodes) > 0 {
        return blk.nodes[0]
    } // if

    return nil
}

func (blk *NodeBlock) scanNext() *Node {
    if blk.current + 1 < len(blk.nodes) {
        blk.current ++
        return blk.nodes[blk.current]
    } // if
    
    return nil
}

func NewNodeBlock() *NodeBlock {
    return &NodeBlock{}
}

type Arc struct {
    head *Node /* node the arc points to */
    next *Arc /* next arc with the same originating node */
    sister *Arc /* reverse arc */
    rCap CapType /* residual capacity */
}

type Node struct {
    first *Arc /* first outcoming arc */
    parent *Arc /* node's parent */
    next *Node /* pointer to the next active node (or to itself if it is the last node in the list) */
    TS int /* timestamp showing when DIST was computed */
    DIST int /* distance to the terminal */
    isSink bool /* flag showing whether the node is in the source or in the sink tree */
    trCap CapType /* if tr_cap > 0 then tr_cap is residual capacity of the arc SOURCE->node
                     otherwise         -tr_cap is residual capacity of the arc node->SINK */
}

var TERMINAL *Arc = &Arc{}
var ORPHAN *Arc = &Arc{}
const INFINITE_D int = 1000000000

type NodePtr struct {
    ptr *Node
    next *NodePtr
}

type Graph struct {
    nodeBlock *NodeBlock
    
    flow CapType /* total flow */
    queueFirst [2]*Node
    queueLast [2]*Node
    orphanFirst *NodePtr
    orphanLast *NodePtr
    TIME int
}

func NewGraph() *Graph {
    return &Graph{nodeBlock: NewNodeBlock()}
}

func (g *Graph) SetTweights(i *Node, capSource, capSink CapType) {
    if capSource < capSink {
        g.flow += capSource
    } else {
        g.flow += capSink
    } // else
    
    i.trCap = capSource - capSink
}

func (g *Graph) AddEdge(from, to *Node, cap, revCap CapType) {
	a, aRev := &Arc{}, &Arc{}
	

	a.sister = aRev
	aRev.sister = a

	a.next = from.first
	from.first = a

	aRev.next = to.first
	to.first = aRev

	a.head = to
	aRev.head = from

	a.rCap = cap
	aRev.rCap = revCap
}

func (g *Graph) Flow() CapType {
    return g.flow
}

func (g *Graph) Maxflow() CapType {
    var i, j, currentNode *Node
    var a *Arc
    var np, npNext *NodePtr

	g.maxflowInit()

    for {
		if i = currentNode; currentNode != nil {
		    i.next = nil
		    if i.parent == nil {
		        i = nil
		    } //  if
		} // if
		
		if i == nil {
		    if i = g.nextActive(); i == nil {
		        break
		    } // if
		} // if

		/* growth */
		if !i.isSink {
			/* grow source tree */
			for a = i.first; a != nil; a = a.next {
			    if a.rCap != 0 {
				    j = a.head
    				if j.parent == nil {
    				    j.isSink = false
    				    j.parent = a.sister
    				    j.TS = i.TS
    				    j.DIST = i.DIST + 1
    				    g.setActive(j)
    				} else if j.isSink {
    				    break
    				} else if j.TS <= i.TS && j.DIST > i.DIST {
    				    /* heuristic - trying to make the distance from j to the source shorter */
    				    j.parent = a.sister
    				    j.TS = i.TS
    				    j.DIST = i.DIST + 1
    				} // else if
    			} // if
			} // for a
        } else {
			/* grow sink tree */
			for a = i.first; a != nil; a = a.next {
			    if a.sister.rCap != 0 {
				    j = a.head
				    if j.parent == nil {
				        j.isSink = true
				        j.parent = a.sister
				        j.TS = i.TS
				        j.DIST = i.DIST + 1
				        g.setActive(j)
				    } else if !j.isSink {
				        a = a.sister
				        break
				    } else if j.TS <= i.TS && j.DIST > i.DIST {
				        /* heuristic - trying to make the distance from j to the sink shorter */
				        j.parent = a.sister
				        j.TS = i.TS
				        j.DIST = i.DIST + 1
				    } // else if
				} // if
            } // for a
	    } // else

		g.TIME ++

		if a != nil {
		    /* set active flag */
			i.next = i
			currentNode = i

			/* augmentation */
			g.augment(a)
			/* augmentation end */

			/* adoption */
			for np = g.orphanFirst; np != nil; np = g.orphanFirst {
				npNext = np.next
				np.next = nil

				for np = g.orphanFirst; np != nil; np = g.orphanFirst {
					//nodeptr_block -> Delete(np);  TODO reuse nodeptr
					g.orphanFirst = np.next
					i = np.ptr
					if g.orphanFirst == nil {
					    g.orphanLast = nil
					} // if
					if i.isSink {
					    g.processSinkOrphan(i)
					} else {
					    g.processSourceOrphan(i)
					} // else
			    } // for np

				g.orphanFirst = npNext
		    } // for np
			/* adoption end */
	    } else {
	        currentNode = nil
	    } // else
    } // for true

	return g.flow
}

func (g *Graph) AddNode() *Node {
    return g.nodeBlock.New()
}

func (g *Graph) IsSource(i *Node) bool {
	return i.parent != nil && !i.isSink
}

func (g *Graph) maxflowInit() {
	g.queueFirst[0], g.queueFirst[1] = nil, nil
	g.queueLast[0], g.queueLast[1] = nil, nil

	for i := g.nodeBlock.scanFirst(); i != nil; i = g.nodeBlock.scanNext() {
		i.next = nil
		i.TS = 0
		if i.trCap > 0 {
			/* i is connected to the source */
			i.isSink = false
			i.parent = TERMINAL
			g.setActive(i)
			i.TS = 0
			i.DIST = 1
	    } else if i.trCap < 0 {
			/* i is connected to the sink */
			i.isSink = true
			i.parent = TERMINAL
			g.setActive(i)
			i.TS = 0
			i.DIST = 1
	    } else {
			i.parent = nil
	    } // else
    } // for i
	g.TIME = 0
}

/*
	Returns the next active node.
	If it is connected to the sink, it stays in the list,
	otherwise it is removed from the list
*/
func (g *Graph) nextActive() *Node {
	var i *Node

	for {
		if i = g.queueFirst[0]; i == nil {
			g.queueFirst[0], i = g.queueFirst[1], g.queueFirst[1]
			g.queueLast[0] = g.queueLast[1]
			g.queueFirst[1] = nil
			g.queueLast[1] = nil
			if i == nil {
			    return nil
			} // if
		} // if

		/* remove it from the active list */
		if i.next == i {
		    g.queueFirst[0], g.queueLast[0] = nil, nil
		} else {
		    g.queueFirst[0] = i.next
		} // else
		i.next = nil

		/* a node in the list is active iff it has a parent */
		if i.parent != nil {
		    return i
		} // if
    } // for true
    
    return nil
}

func (g *Graph) setActive(i *Node) {
	if i.next == nil {
		/* it's not in the list yet */
		if g.queueLast[1] != nil {
		    g.queueLast[1].next = i
		} else {
		    g.queueFirst[1] = i
		} // else
		g.queueLast[1] = i
		i.next = i
	} // if
}

func (g *Graph) augment(middleArc *Arc) {
	var i *Node
	var a *Arc
	var bottleNeck CapType
	var np *NodePtr


	/* 1. Finding bottleneck capacity */
	/* 1a - the source tree */
	bottleNeck = middleArc.rCap
	for i = middleArc.sister.head; true; i = a.head {
		a = i.parent
		if a == TERMINAL {
		    break
		} // if
		if bottleNeck > a.sister.rCap {
		    bottleNeck = a.sister.rCap
		} // if
	} // for i
	if bottleNeck > i.trCap {
	    bottleNeck = i.trCap
	} // if
	/* 1b - the sink tree */
	for i = middleArc.head; true; i = a.head {
		a = i.parent
		if a == TERMINAL {
		    break
		} // if
		if bottleNeck > a.rCap {
		    bottleNeck = a.rCap
		} // if
	} // for i
	if bottleNeck > -i.trCap {
	    bottleNeck = - i.trCap
	} // if


	/* 2. Augmenting */
	/* 2a - the source tree */
	middleArc.sister.rCap += bottleNeck
	middleArc.rCap -= bottleNeck
	for i = middleArc.sister.head; true; i = a.head {
		a = i.parent
		if a == TERMINAL {
		    break
		} // if
		a.rCap += bottleNeck
		a.sister.rCap -= bottleNeck
		if a.sister.rCap == 0 {
			/* add i to the adoption list */
			i.parent = ORPHAN
			np = &NodePtr{}
			np.ptr = i
			np.next = g.orphanFirst
			g.orphanFirst = np
		} // if
	} // for i
	i.trCap -= bottleNeck
	if i.trCap == 0 {
		/* add i to the adoption list */
        i.parent = ORPHAN
        np = &NodePtr{}
        np.ptr = i
        np.next = g.orphanFirst
        g.orphanFirst = np
	} // if
	/* 2b - the sink tree */
	for i = middleArc.head; true; i = a.head {
		a = i.parent
		if a == TERMINAL {
		    break
		} // if
		a.sister.rCap += bottleNeck
		a.rCap -= bottleNeck
		if a.rCap == 0 {
			/* add i to the adoption list */
			i.parent = ORPHAN
			np = &NodePtr{}
			np.ptr = i
			np.next = g.orphanFirst
			g.orphanFirst = np
		} // if
	} // for i
	i.trCap += bottleNeck
	if i.trCap == 0 {
		/* add i to the adoption list */
        i.parent = ORPHAN
        np = &NodePtr{}
        np.ptr = i
        np.next = g.orphanFirst
        g.orphanFirst = np
	} // if

	g.flow += bottleNeck
}

func (g *Graph) processSinkOrphan(i *Node) {
	var a0Min *Arc
	var dMin int = INFINITE_D

	/* trying to find a new parent */
	for a0 := i.first; a0 != nil; a0 = a0.next {
	    if a0.rCap != 0 {
		    j := a0.head
		    if a := j.parent; j.isSink && a != nil {
			/* checking the origin of j */
			//d = 0;
			    var d int = 0
			    for true {
				    if j.TS == g.TIME {
				        d += j.DIST
				        break
				    } // if
				    a = j.parent
				    d ++
				    if a == TERMINAL {
					    j.TS = g.TIME
					    j.DIST = 1
					    break
					} // if
    				if a == ORPHAN {
    				    d = INFINITE_D
    				    break
    				} // if
				    j = a.head
			    } // for true
    			/* j originates from the sink - done */
			    if d < INFINITE_D {
				    if d < dMin {
				        a0Min = a0
				        dMin = d
				    } // if
                    /* set marks along the path */
    				for j := a0.head; j.TS != g.TIME; j = j.parent.head {
    				    j.TS = g.TIME
    				    j.DIST = d
    				    d --
    				} // for j
			    } // if
		    } // if
	    } // if
	} // for a0

	if i.parent = a0Min; i.parent != nil {
		i.TS = g.TIME
		i.DIST = dMin + 1
    } else {
		/* no parent is found */
		i.TS = 0

		/* process neighbors */
		for a0 := i.first; a0 != nil; a0 = a0.next {
			j := a0.head
			if a := j.parent; j.isSink && a != nil {
				if a0.rCap != 0 {
				    g.setActive(j)
				} // if
				if a != TERMINAL && a != ORPHAN && a.head == i {
					/* add j to the adoption list */
					i.parent = ORPHAN
					np := &NodePtr{}
					np.ptr = j
					if g.orphanLast != nil {
					    g.orphanLast.next = np
					} else {
					    g.orphanFirst = np
					} // else
					g.orphanLast = np
					np.next = nil
				} // i f
			} // if
		} // for a0
	} // else
}

func (g *Graph) processSourceOrphan(i *Node) {
	var a0Min *Arc
	var dMin int = INFINITE_D

	/* trying to find a new parent */
	for a0 := i.first; a0 != nil; a0 = a0.next {
    	if a0.sister.rCap != 0 {
    		j := a0.head
    		if a := j.parent; j.isSink && a != nil {
    			/* checking the origin of j */
    			var d int = 0
    			for true {
    				if j.TS == g.TIME {
    				    d += j.DIST
    				    break
    				} // if
    				a = j.parent
    				d ++
    				if a == TERMINAL {
    				    j.TS = g.TIME
    				    j.DIST = 1
    				    break
    				} // if
    				if a == ORPHAN {
    				    d = INFINITE_D
    				    break
    				} // if
    				j = a.head
    			} // for true
                /* j originates from the source - done */
    			if d < INFINITE_D {
    				if d < dMin {
    				    a0Min = a0
    				    dMin = d
    				} // if
    				/* set marks along the path */
    				for j := a0.head; j.TS != g.TIME; j = j.parent.head {
    				    j.TS = g.TIME
    				    j.DIST = d
    				    d --
    				} // for j
    			} // if
  		    } // if
    	} // if
    } // for a0

	if i.parent = a0Min; i.parent != nil {
		i.TS = g.TIME
		i.DIST = dMin + 1
    } else {
		/* no parent is found */
		i.TS = 0

		/* process neighbors */
		for a0 := i.first; a0 != nil; a0 = a0.next {
			j := a0.head
			if a := j.parent; !j.isSink && a != nil {
				if a0.sister.rCap != 0 {
				    g.setActive(j)
				} // if
				if a != TERMINAL && a != ORPHAN && a.head == i {
					/* add j to the adoption list */
					j.parent = ORPHAN
					np := &NodePtr{}
					np.ptr = j
					if g.orphanLast != nil {
					    g.orphanLast.next = np
					} else {
					    g.orphanFirst = np
					} // else
					g.orphanLast = np
					np.next = nil
				} // if
			} //  if
		} // for a0
	} // else
}
