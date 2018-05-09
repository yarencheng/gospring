package dependency

type Graph struct {
	nodes map[string]*Node
}
type Node struct {
	parents map[string]*Node
	value   string
}

func NewGraph() *Graph {
	return &Graph{
		nodes: make(map[string]*Node),
	}
}

func NewNode(value string) *Node {
	return &Node{
		parents: make(map[string]*Node),
		value:   value,
	}
}

func (g *Graph) AddDependency(parent string, child string) (ok bool) {

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

func (g *Graph) isLoop(parent string, child string) bool {

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
