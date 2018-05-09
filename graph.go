package gospring

type graph struct {
	nodes map[string]*node
}
type node struct {
	parents map[string]*node
	value   string
}

func newGraph() *graph {
	return &graph{
		nodes: make(map[string]*node),
	}
}

func NewNode(value string) *node {
	return &node{
		parents: make(map[string]*node),
		value:   value,
	}
}

func (g *graph) AddDependency(parent string, child string) (ok bool) {

	if g.isLoop(parent, child) {
		return false
	}

	pNode, pResent := g.nodes[parent]
	if !pResent {
		pNode = NewNode(parent)
		g.nodes[parent] = pNode
	}

	cNode, cResent := g.nodes[child]
	if !cResent {
		cNode = NewNode(child)
		g.nodes[child] = cNode
	}

	cNode.parents[parent] = pNode

	return true
}

func (g *graph) isLoop(parent string, child string) bool {

	if parent == child {
		return true
	}

	pNode, pResent := g.nodes[parent]
	if !pResent {
		return false
	}

	_, cResent := g.nodes[child]
	if !cResent {
		return false
	}

	for _, p := range pNode.parents {
		if g.isLoop(p.value, child) {
			return true
		}
	}

	return false
}
