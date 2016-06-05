package algo_test

import (
	"reflect"
	"testing"

	"github.com/octo47/gomisc/algo"
)

func TestSimplePath(t *testing.T) {
	graph := algo.NewGraph()
	v1 := graph.AddVertex()
	v2 := graph.AddVertex()
	v3 := graph.AddVertex()
	graph.AddEdge(v1, v2, 3, false)
	graph.AddEdge(v1, v3, 2, false)
	graph.AddEdge(v2, v3, -2, false)

	{
		path, err := graph.Dijkstra(v1)
		if err != nil {
			t.Error(err)
		}

		if pathTo3 := path.BuildPath(v3); !reflect.DeepEqual(pathTo3, []int{2, 1, 0}) {
			t.Error("Wrong path calculated", pathTo3, path)
		}

		if cost := path.PathCost(v3); cost != 1 {
			t.Error("Wrong path cost calculated", path)
		}
	}

	graph.DelEdge(v2, v3)
	{
		path, err := graph.Dijkstra(v1)
		if err != nil {
			t.Error(err)
		}

		if pathTo3 := path.BuildPath(v3); !reflect.DeepEqual(pathTo3, []int{2, 0}) {
			t.Error("Wrong path calculated", pathTo3, path)
		}
	}
}
