package algo

import (
	"math"
)

const Undef int = -1
const UndefDist int = math.MaxInt32

type edge struct {
	target int
	cost   int
}

type Graph struct {
	edges [][]edge
}

type Path struct {
	source int
	dist   []int
	prev   []int
}

func NewGraph() *Graph {
	s := &Graph{
		edges: make([][]edge, 0),
	}
	return s
}

func (p *Path) BuildPath(target int) []int {
	stack := NewStack()
	u := Undef
	for u = target; u != Undef; {
		stack.Push(u)
		u = p.prev[u]
	}
	return stack.AsSlice()
}

func (p *Path) PathCost(target int) int {
	return p.dist[target]
}

// add vertex to graph and return its index
func (g *Graph) AddVertex() int {
	g.edges = append(g.edges, make([]edge, 0))
	return len(g.edges) - 1
}

func (g *Graph) AddVertexes(count int) {
	for i := 0; i < count; i++ {
		g.edges = append(g.edges, make([]edge, 0))
	}
}

func (g *Graph) DelEdge(from int, to int) {
	for idx, edge := range g.edges[from] {
		if edge.target == to {
			a := g.edges[from]
			g.edges[from] = append(a[:idx], a[idx+1:]...)
		}
	}
}

func (g *Graph) AddEdge(vertex1 int, vertex2 int, cost int, bidir bool) {
	g.edges[vertex1] = append(g.edges[vertex1], edge{target: vertex2, cost: cost})
	if bidir {
		g.edges[vertex2] = append(g.edges[vertex2], edge{target: vertex1, cost: cost})
	}
}

func (g *Graph) Dijkstra(source int) (*Path, error) {
	n := len(g.edges)
	dist := make([]int, n, n)
	prev := make([]int, n, n)
	visited := make([]bool, n, n)
	for i := 0; i < n; i++ {
		dist[i] = math.MaxInt32
		prev[i] = Undef
	}
	dist[source] = 0

	for i := 0; i < n; i++ {
		u := Undef
		for j := 0; j < n; j++ {
			if visited[j] {
				continue
			}
			if u == Undef || dist[j] < dist[u] {
				u = j
			}
		}
		visited[u] = true
		for ni := 0; ni < len(g.edges[u]); ni++ {
			e := g.edges[u][ni]
			v := e.target
			alt := dist[u] + e.cost
			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u
			}
		}
	}
	return &Path{source: source, dist: dist, prev: prev}, nil
}
