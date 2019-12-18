package dijkstra_test

import (
	"testing"

	algo "github.com/octo47/gomisc/algo/dijkstra"
)

func TestSimpleStack(t *testing.T) {
	st := algo.NewStack()
	st.Push(1)
	st.Push(2)
	st.Push(3)
	if v, _ := st.Pop(); v != 3 {
		t.Fail()
	}
	if v, _ := st.Pop(); v != 2 {
		t.Fail()
	}
	if v, _ := st.Pop(); v != 1 {
		t.Fail()
	}
}
