package maxflow

import (
    "fmt"
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
	// flow += (cap_source < cap_sink) ? cap_source : cap_sink;
	// ((node*)i) -> tr_cap = cap_source - cap_sink;
    if capSource < capSink {
        g.flow += capSource
    } else {
        g.flow += capSink
    } // else
    
    i.trCap = capSource - capSink
}

func (g *Graph) AddEdge(from, to *Node, cap, revCap CapType) {
	//arc *a, *a_rev;

	//a = arc_block -> New(2);
	//a_rev = a + 1;
	a, aRev := &Arc{}, &Arc{}
	

	//a -> sister = a_rev;
	//a_rev -> sister = a;
	a.sister = aRev
	aRev.sister = a
	//a -> next = ((node*)from) -> first;
	//((node*)from) -> first = a;
	a.next = from.first
	from.first = a
	//a_rev -> next = ((node*)to) -> first;
	//((node*)to) -> first = a_rev;
	aRev.next = to.first
	to.first = aRev
	//a -> head = (node*)to;
	//a_rev -> head = (node*)from;
	a.head = to
	aRev.head = from
	//a -> r_cap = cap;
	//a_rev -> r_cap = rev_cap;
	a.rCap = cap
	aRev.rCap = revCap
}

func (g *Graph) maxflowInit() {
	//node *i;

	//queue_first[0] = queue_last[0] = NULL;
	//queue_first[1] = queue_last[1] = NULL;
	//orphan_first = NULL;
	g.queueFirst[0], g.queueFirst[1] = nil, nil
	g.queueLast[0], g.queueLast[1] = nil, nil

	//for (i=node_block->ScanFirst(); i; i=node_block->ScanNext())
	//{
	for i := g.nodeBlock.scanFirst(); i != nil; i = g.nodeBlock.scanNext() {
		//i -> next = NULL;
		//i -> TS = 0;
		i.next = nil
		i.TS = 0
		//if (i->tr_cap > 0)
		//{
		if i.trCap > 0 {
			/* i is connected to the source */
			//i -> is_sink = 0;
			//i -> parent = TERMINAL;
			//set_active(i);
			//i -> TS = 0;
			//i -> DIST = 1;
			i.isSink = false
			i.parent = TERMINAL
			g.setActive(i)
			i.TS = 0
			i.DIST = 1
		//}
		//else if (i->tr_cap < 0)
	//	{
	    } else if i.trCap < 0 {
			/* i is connected to the sink */
			//i -> is_sink = 1;
			//i -> parent = TERMINAL;
			//set_active(i);
			//i -> TS = 0;
			//i -> DIST = 1;
			i.isSink = true
			i.parent = TERMINAL
			g.setActive(i)
			i.TS = 0
			i.DIST = 1
		//}
		//else
	    //	{
	    } else {
			//i -> parent = NULL;
			i.parent = nil
		//}
	    } // else
	//}
    } // for i
	//TIME = 0;
	g.TIME = 0
}

/*
	Returns the next active node.
	If it is connected to the sink, it stays in the list,
	otherwise it is removed from the list
*/
func (g *Graph) nextActive() *Node {
	//node *i;
	var i *Node

	//while ( 1 )
	//{
	for {
		//if (!(i=queue_first[0]))
		//{
		if i = g.queueFirst[0]; i == nil {
			//queue_first[0] = i = queue_first[1];
			//queue_last[0]  = queue_last[1];
			//queue_first[1] = NULL;
			//queue_last[1]  = NULL;
			//if (!i) return NULL;
			g.queueFirst[0], i = g.queueFirst[1], g.queueFirst[1]
			g.queueLast[0] = g.queueLast[1]
			g.queueFirst[1] = nil
			g.queueLast[1] = nil
			if i == nil {
			    return nil
			} // if
		//}
		} // if

		/* remove it from the active list */
		//if (i->next == i) queue_first[0] = queue_last[0] = NULL;
		//else              queue_first[0] = i -> next;
		if i.next == i {
		    g.queueFirst[0], g.queueLast[0] = nil, nil
		} else {
		    g.queueFirst[0] = i.next
		} // else
		//i -> next = NULL;
		i.next = nil

		/* a node in the list is active iff it has a parent */
		//if (i->parent) return i;
		if i.parent != nil {
		    return i
		} // if
	//}
    } // for true
    
    return nil
}

func (g *Graph) setActive(i *Node) {
	//if (!i->next)
	//{
	if i.next == nil {
		/* it's not in the list yet */
		//if (queue_last[1]) queue_last[1] -> next = i;
		//else               queue_first[1]        = i;
		if g.queueLast[1] != nil {
		    g.queueLast[1].next = i
		} else {
		    g.queueFirst[1] = i
		} // else
		//queue_last[1] = i;
		//i -> next = i;
		g.queueLast[1] = i
		i.next = i
	//}
	} // if
}

func (g *Graph) augment(middleArc *Arc) {
	//node *i;
	//arc *a;
	//captype bottleneck;
	//nodeptr *np;
	var i *Node
	var a *Arc
	var bottleNeck CapType
	var np *NodePtr


	/* 1. Finding bottleneck capacity */
	/* 1a - the source tree */
	//bottleneck = middle_arc -> r_cap;
	bottleNeck = middleArc.rCap
	//for (i=middle_arc->sister->head; ; i=a->head)
	//{
	for i = middleArc.sister.head; true; i = a.head {
		//a = i -> parent;
		a = i.parent
		//if (a == TERMINAL) break;
		if a == TERMINAL {
		    break
		} // if
		//if (bottleneck > a->sister->r_cap) bottleneck = a -> sister -> r_cap;
		if bottleNeck > a.sister.rCap {
		    bottleNeck = a.sister.rCap
		} // if
	//}
	} // for i
	//if (bottleneck > i->tr_cap) bottleneck = i -> tr_cap;
	if bottleNeck > i.trCap {
	    bottleNeck = i.trCap
	} // if
	/* 1b - the sink tree */
	//for (i=middle_arc->head; ; i=a->head)
	//{
	for i = middleArc.head; true; i = a.head {
		//a = i -> parent;
		a = i.parent
		//if (a == TERMINAL) break;
		if a == TERMINAL {
		    break
		} // if
		//if (bottleneck > a->r_cap) bottleneck = a -> r_cap;
		if bottleNeck > a.rCap {
		    bottleNeck = a.rCap
		} // if
	//}
	} // for i
	//if (bottleneck > - i->tr_cap) bottleneck = - i -> tr_cap;
	if bottleNeck > -i.trCap {
	    bottleNeck = - i.trCap
	} // if


	/* 2. Augmenting */
	/* 2a - the source tree */
	//middle_arc -> sister -> r_cap += bottleneck;
	//middle_arc -> r_cap -= bottleneck;
	middleArc.sister.rCap += bottleNeck
	middleArc.rCap -= bottleNeck
	//for (i=middle_arc->sister->head; ; i=a->head)
	//{
	for i = middleArc.sister.head; true; i = a.head {
		//a = i -> parent;
		a = i.parent
		//if (a == TERMINAL) break;
		if a == TERMINAL {
		    break
		} // if
		//a -> r_cap += bottleneck;
		//a -> sister -> r_cap -= bottleneck;
		a.rCap += bottleNeck
		a.sister.rCap -= bottleNeck
		//if (!a->sister->r_cap)
		//{
		if a.sister.rCap == 0 {
			/* add i to the adoption list */
			//i -> parent = ORPHAN;
			//np = nodeptr_block -> New();
			//np -> ptr = i;
			//np -> next = orphan_first;
			//orphan_first = np;
			i.parent = ORPHAN
			np = &NodePtr{}
			np.ptr = i
			np.next = g.orphanFirst
			g.orphanFirst = np
		 //}
		} // if
	//}
	} // for i
	//i -> tr_cap -= bottleneck;
	i.trCap = bottleNeck
	//if (!i->tr_cap)
	//{
	if i.trCap == 0 {
		/* add i to the adoption list */
		//i -> parent = ORPHAN;
		//np = nodeptr_block -> New();
		//np -> ptr = i;
		//np -> next = orphan_first;
		//orphan_first = np;
        i.parent = ORPHAN
        np = &NodePtr{}
        np.ptr = i
        np.next = g.orphanFirst
        g.orphanFirst = np
	//}
	} // if
	/* 2b - the sink tree */
	//for (i=middle_arc->head; ; i=a->head)
	//{
	for i = middleArc.head; true; i = a.head {
		//a = i -> parent;
		a = i.parent
		//if (a == TERMINAL) break;
		if a == TERMINAL {
		    break
		} // if
		//a -> sister -> r_cap += bottleneck;
		//a -> r_cap -= bottleneck;
		a.sister.rCap += bottleNeck
		a.rCap -= bottleNeck
		//if (!a->r_cap)
		//{
		if a.rCap == 0 {
			/* add i to the adoption list */
			//i -> parent = ORPHAN;
			//np = nodeptr_block -> New();
			//np -> ptr = i;
			//np -> next = orphan_first;
			//orphan_first = np;
			i.parent = ORPHAN
			np = &NodePtr{}
			np.ptr = i
			np.next = g.orphanFirst
			g.orphanFirst = np
		//}
		} // if
	//}
	} // for i
	//i -> tr_cap += bottleneck;
	i.trCap += bottleNeck
	//if (!i->tr_cap)
	//{
	if i.trCap == 0 {
		/* add i to the adoption list */
		//i -> parent = ORPHAN;
		//np = nodeptr_block -> New();
		//np -> ptr = i;
		//np -> next = orphan_first;
		//orphan_first = np;
        i.parent = ORPHAN
        np = &NodePtr{}
        np.ptr = i
        np.next = g.orphanFirst
        g.orphanFirst = np
	//}
	} // if

	//flow += bottleneck;
	g.flow += bottleNeck
}

func (g *Graph) processSinkOrphan(i *Node) {
	//node *j;
	//arc *a0, *a0_min = NULL, *a;
	//nodeptr *np;
	//int d, d_min = INFINITE_D;
	var a0Min *Arc
	var dMin int = INFINITE_D

	/* trying to find a new parent */
	//for (a0=i->first; a0; a0=a0->next)
	for a0 := i.first; a0 != nil; a0 = a0.next {
	//if (a0->r_cap)
	//{
	    if a0.rCap != 0 {
		//j = a0 -> head;
		    j := a0.head
		//if (j->is_sink && (a=j->parent))
		//{
		    if a := j.parent; j.isSink && a != nil {
			/* checking the origin of j */
			//d = 0;
			    d := int(0)
			//while ( 1 )
			//{
			    for true {
				//if (j->TS == TIME)
				//{
				//	d += j -> DIST;
				//	break;
				//}
				    if j.TS == g.TIME {
				        d += j.DIST
				        break
				    } // if
				//a = j -> parent;
				//d ++;
				    a = j.parent
				    d ++
				//if (a==TERMINAL)
				//{
					//j -> TS = TIME;
					//j -> DIST = 1;
					//break;
				//}
				    if a == TERMINAL {
					    j.TS = g.TIME
					    j.DIST = 1
					    break
					} // if
				//if (a==ORPHAN) { d = INFINITE_D; break; }
				if a == ORPHAN {
				    d = INFINITE_D
				    break
				} // if
				//j = a -> head;
				    j = a.head
			//}
			    } // for true
			//if (d<INFINITE_D) /* j originates from the sink - done */
			//{
			    if d < INFINITE_D {
				//if (d<d_min)
				//{
				//	a0_min = a0;
				//	d_min = d;
				//}
				    if d < dMin {
				        a0Min = a0
				        dMin = d
				    } // if
				/* set marks along the path */
				//for (j=a0->head; j->TS!=TIME; j=j->parent->head)
				//{
					//j -> TS = TIME;
					//j -> DIST = d --;
				//}
				for j := a0.head; j.TS != g.TIME; j = j.parent.head {
				    j.TS = g.TIME
				    j.DIST = d
				    d --
				} // for j
			//}
			    } // if
		//}
		    } // if
	//}
	    } // if
	} // for a0

	//if (i->parent = a0_min)
	//{
	if i.parent = a0Min; i.parent != nil {
		//i -> TS = TIME;
		//i -> DIST = d_min + 1;
		i.TS = g.TIME
		i.DIST = dMin + 1
	//}
	//else
	//{
    } else {
		/* no parent is found */
		//i -> TS = 0;
		i.TS = 0

		/* process neighbors */
		//for (a0=i->first; a0; a0=a0->next)
		//{
		for a0 := i.first; a0 != nil; a0 = a0.next {
			//j = a0 -> head;
			j := a0.head
			//if (j->is_sink && (a=j->parent))
			//{
			if a := j.parent; j.isSink && a != nil {
				//if (a0->r_cap) set_active(j);
				if a0.rCap != 0 {
				    g.setActive(j)
				} // if
				//if (a!=TERMINAL && a!=ORPHAN && a->head==i)
				//{
				if a != TERMINAL && a != ORPHAN && a.head == i {
					/* add j to the adoption list */
					//j -> parent = ORPHAN;
					//np = nodeptr_block -> New();
					//np -> ptr = j;
					i.parent = ORPHAN
					np := &NodePtr{}
					np.ptr = j
					//if (orphan_last) orphan_last -> next = np;
					//else             orphan_first        = np;
					if g.orphanLast != nil {
					    g.orphanLast.next = np
					} else {
					    g.orphanFirst = np
					} // else
					//orphan_last = np;
					//np -> next = NULL;
					g.orphanLast = np
					np.next = nil
				//}
				} // i f
			//}
			} // if
		 //}
		} // for a0
	//}
	} // else
}

func (g *Graph) processSourceOrphan(i *Node) {
	//node *j;
	//arc *a0, *a0_min = NULL, *a;
	//nodeptr *np;
	//int d, d_min = INFINITE_D;
	var a0Min *Arc
	var dMin int = INFINITE_D

	/* trying to find a new parent */
	//for (a0=i->first; a0; a0=a0->next)
	for a0 := i.first; a0 != nil; a0 = a0.next {
    	//if (a0->sister->r_cap)
    	//{
    	if a0.sister.rCap != 0 {
    		//j = a0 -> head;
    		j := a0.head
    		//if (!j->is_sink && (a=j->parent))
    		//{
    		if a := j.parent; j.isSink && a != nil {
    			/* checking the origin of j */
    			//d = 0;
    			var d int = 0
    			//while ( 1 )
    			//{
    			for true {
    				//if (j->TS == TIME)
    				//{
    				//	d += j -> DIST;
    				//	break;
    				//}
    				if j.TS == g.TIME {
    				    d += j.DIST
    				    break
    				} // if
    				//a = j -> parent;
    				//d ++;
    				a = j.parent
    				d ++
    				//if (a==TERMINAL)
    				//{
    				//	j -> TS = TIME;
    				//	j -> DIST = 1;
    				//	break;
    				//}
    				if a == TERMINAL {
    				    j.TS = g.TIME
    				    j.DIST = 1
    				    break
    				} // if
    				//if (a==ORPHAN) { d = INFINITE_D; break; }
    				if a == ORPHAN {
    				    d = INFINITE_D
    				    break
    				} // if
    				//j = a -> head;
    				j = a.head
    			//}
    			} // for true
    			//if (d<INFINITE_D) /* j originates from the source - done */
    			//{
    			if d < INFINITE_D {
    				//if (d<d_min)
    				//{
    				//	a0_min = a0;
    				//	d_min = d;
    				//}
    				if d < dMin {
    				    a0Min = a0
    				    dMin = d
    				} // if
    				/* set marks along the path */
    				//for (j=a0->head; j->TS!=TIME; j=j->parent->head)
    				//{
    				//	j -> TS = TIME;
    				//	j -> DIST = d --;
    				//}
    				for j := a0.head; j.TS != g.TIME; j = j.parent.head {
    				    j.TS = g.TIME
    				    j.DIST = d
    				    d --
    				} // for j
    			//}
    			} // if
  		    //}
  		    } // if
    	//}
    	} // if
    } // for a0

	//if (i->parent = a0_min)
	//{
	if i.parent = a0Min; i.parent != nil {
		//i -> TS = TIME;
		//i -> DIST = d_min + 1;
		i.TS = g.TIME
		i.DIST = dMin + 1
	//}
	//else
	//{
    } else {
		/* no parent is found */
		//i -> TS = 0;
		i.TS = 0

		/* process neighbors */
		//for (a0=i->first; a0; a0=a0->next)
		//{
		for a0 := i.first; a0 != nil; a0 = a0.next {
			//j = a0 -> head;
			j := a0.head
			//if (!j->is_sink && (a=j->parent))
			//{
			if a := j.parent; !j.isSink && a != nil {
				//if (a0->sister->r_cap) set_active(j);
				if a0.sister.rCap != 0 {
				    g.setActive(j)
				} // if
				//if (a!=TERMINAL && a!=ORPHAN && a->head==i)
				//{
				if a != TERMINAL && a != ORPHAN && a.head == i {
					/* add j to the adoption list */
					//j -> parent = ORPHAN;
					//np = nodeptr_block -> New();
					//np -> ptr = j;
					j.parent = ORPHAN
					np := &NodePtr{}
					np.ptr = j
					//if (orphan_last) orphan_last -> next = np;
					//else             orphan_first        = np;
					if g.orphanLast != nil {
					    g.orphanLast.next = np
					} else {
					    g.orphanFirst = np
					} // else
					//orphan_last = np;
					//np -> next = NULL;
					g.orphanLast = np
					np.next = nil
				//}
				} // if
			//}
			} //  if
		//}
		} // for a0
	//}
	} // else
}

func (g *Graph) Maxflow() CapType {
    /*
	node *i, *j, *current_node = NULL;
	arc *a;
	nodeptr *np, *np_next;
	*/
    var i, j, currentNode *Node
    var a *Arc
    var np, npNext *NodePtr

	//maxflow_init();
	g.maxflowInit()
	//
	//nodeptr_block = new DBlock<nodeptr>(NODEPTR_BLOCK_SIZE, error_function);

	//while ( 1 )
	//{
    for {
		//if (i=current_node)
		//{
		//	i -> next = NULL; /* remove active flag */
		//	if (!i->parent) i = NULL;
		//}
		i = currentNode
		if currentNode != nil {
		    i.next = nil
		    if i.parent == nil {
		        i = nil
		    } //  if
		} // if
		
		//if (!i)
		//{
		//	if (!(i = next_active())) break;
		//}
		if i == nil {
		    i = g.nextActive()
		    if i == nil {
		        break
		    } // if
		} // if

		/* growth */
		//if (!i->is_sink)
		//{
		if !i.isSink {
			/* grow source tree */
			//for (a=i->first; a; a=a->next)
			for a = i.first; a != nil; a = a.next {
			//if (a->r_cap)
			//{
			    if a.rCap != 0 {
				//j = a -> head;
				    j = a.head
				//if (!j->parent)
				//{
				//	j -> is_sink = 0;
				//	j -> parent = a -> sister;
				//	j -> TS = i -> TS;
				//	j -> DIST = i -> DIST + 1;
				//	set_active(j);
				//}
				//else if (j->is_sink) break;
				//else if (j->TS <= i->TS &&
				//         j->DIST > i->DIST)
				//{
					/* heuristic - trying to make the distance from j to the source shorter */
				//	j -> parent = a -> sister;
				//	j -> TS = i -> TS;
				//	j -> DIST = i -> DIST + 1;
				//}
    				if j.parent == nil {
    				    j.isSink = false
    				    j.parent = a.sister
    				    j.TS = i.TS
    				    j.DIST = i.DIST + 1
    				    g.setActive(j)
    				} else if j.isSink {
    				    break
    				} else if j.TS <= i.TS && j.DIST > i.DIST {
    				    j.parent = a.sister
    				    j.TS = i.TS
    				    j.DIST = i.DIST + 1
    				} // else if
    			} // if
			//}
			} // for a
		//}
		//else
		//{
        } else {
			/* grow sink tree */
			//for (a=i->first; a; a=a->next)
			for a = i.first; a != nil; a = a.next {
			//if (a->sister->r_cap)
			//{
			    if a.sister.rCap != 0 {
				//j = a -> head;
				    j = a.head
				//if (!j->parent)
				//{
				//	j -> is_sink = 1;
				//	j -> parent = a -> sister;
				//	j -> TS = i -> TS;
				//	j -> DIST = i -> DIST + 1;
				//	set_active(j);
				//}
				//else if (!j->is_sink) { a = a -> sister; break; }
				//else if (j->TS <= i->TS &&
				//         j->DIST > i->DIST)
				//{
					/* heuristic - trying to make the distance from j to the sink shorter */
				//	j -> parent = a -> sister;
				//	j -> TS = i -> TS;
				//	j -> DIST = i -> DIST + 1;
				//}
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
				        j.parent = a.sister
				        j.TS = i.TS
				        j.DIST = i.DIST + 1
				    } // else if
				} // if
//			}
            } // for a
		//}
	    } // else

		//TIME ++;
		g.TIME ++

		//if (a)
		//{
		if a != nil {
			//i -> next = i; /* set active flag */
			//current_node = i;
			i.next = i
			currentNode = i

			/* augmentation */
			//augment(a);
			g.augment(a)
			/* augmentation end */

			/* adoption */
			//while (np=orphan_first)
			//{
			for np = g.orphanFirst; np != nil; np = g.orphanFirst {
				//np_next = np -> next;
				//np -> next = NULL;
				npNext = np.next
				np.next = nil

				//while (np=orphan_first)
				//{
				for np = g.orphanFirst; np != nil; np = g.orphanFirst {
					//orphan_first = np -> next;
					//i = np -> ptr;
					//nodeptr_block -> Delete(np);
					g.orphanFirst = np.next
					i = np.ptr
					//if (!orphan_first) orphan_last = NULL;
					if g.orphanFirst == nil {
					    g.orphanLast = nil
					} // if
					//if (i->is_sink) process_sink_orphan(i);
					//else            process_source_orphan(i);
					if i.isSink {
					    g.processSinkOrphan(i)
					} else {
					    g.processSourceOrphan(i)
					} // else
				//}
			    } // for np

				//orphan_first = np_next;
				g.orphanFirst = npNext
			//}
		    } // for np
			/* adoption end */
		//}
		//else current_node = NULL;
	    } else {
	        currentNode = nil
	    } // else
	//}
    } // for true
	//delete nodeptr_block;

	//return flow;
	return g.flow
}

func (g *Graph) AddNode() *Node {
    return g.nodeBlock.New()
}

func (g *Graph) IsSource(i *Node) bool {
	//if (((node*)i)->parent && !((node*)i)->is_sink) return SOURCE;
	//return SINK;
	return i.parent != nil && !i.isSink
}

func main() {
    g := NewGraph()
    
//    nodes := [2]*Node{g.AddNode(), g.AddNode()}
    var nodes [4]*Node
    
    for i := range(nodes) {
        nodes[i] = g.AddNode()
    } // for i
    
//    g.SetTweights(nodes[0], 1, 5)
//    g.SetTweights(nodes[1], 2, 6)
//    g.AddEdge(nodes[0], nodes[1], 3, 4)

    g.SetTweights(nodes[0], 3, 0)
    g.SetTweights(nodes[1], 3, 0)
    g.SetTweights(nodes[2], 0, 2)
    g.SetTweights(nodes[3], 0, 3)
    
    g.AddEdge(nodes[0], nodes[1], 2, 0)
    g.AddEdge(nodes[0], nodes[2], 3, 0)
    g.AddEdge(nodes[1], nodes[3], 2, 0)
    g.AddEdge(nodes[2], nodes[3], 4, 0)

	flow := g.Maxflow();
    
    fmt.Println(nodes)
    fmt.Println("Flow =", flow)
    fmt.Println("Minimum cut:")
    for i := range(nodes) {
        if g.IsSource(nodes[i]) {
            fmt.Printf("node%d is in the SOURCE set\n", i)
        } else {
            fmt.Printf("node%d is in the SINK set\n", i)
        } // else
    } // for i
/*    
	Graph::node_id nodes[2];
	Graph *g = new Graph();

	nodes[0] = g -> add_node();
	nodes[1] = g -> add_node();
	g -> set_tweights(nodes[0], 1, 5);
	g -> set_tweights(nodes[1], 2, 6);
	g -> add_edge(nodes[0], nodes[1], 3, 4);

	Graph::flowtype flow = g -> maxflow();

	printf("Flow = %d\n", flow);
	printf("Minimum cut:\n");
	if (g->what_segment(nodes[0]) == Graph::SOURCE)
		printf("node0 is in the SOURCE set\n");
	else
		printf("node0 is in the SINK set\n");
	if (g->what_segment(nodes[1]) == Graph::SOURCE)
		printf("node1 is in the SOURCE set\n");
	else
		printf("node1 is in the SINK set\n");

	delete g;
*/	
}