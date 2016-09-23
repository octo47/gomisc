package mysort

import (
	"sort"
	"testing"
)

var ints = [...]int{
	74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586,
}

func TestMyIntSort(t *testing.T) {
	data := ints
	a := data[0:]
	MyIntSort(a)
	if !sort.IsSorted(sort.IntSlice(a)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestMyIntSortFunc(t *testing.T) {
	data := ints
	a := data[0:]
	MyIntSortFunc(a)
	if !sort.IsSorted(sort.IntSlice(a)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestMyIntSortStruct(t *testing.T) {
	data := ints
	a := data[0:]
	MyIntSortStruct(IntSlice{a})
	if !sort.IsSorted(sort.IntSlice(a)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestMyIntSortPtr(t *testing.T) {
	data := ints
	a := data[0:]
	MyIntSortPtr(&IntSlicePtr{a})
	if !sort.IsSorted(sort.IntSlice(a)) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func BenchmarkSortInt1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		d := sort.IntSlice(data)
		b.StartTimer()
		sort.Sort(d)
		b.StopTimer()
	}
}

func BenchmarkMyIntSort1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		b.StartTimer()
		MyIntSort(data)
		b.StopTimer()
	}
}

func BenchmarkMyIntSortFunc1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		b.StartTimer()
		MyIntSortFunc(data)
		b.StopTimer()
	}
}

func BenchmarkMyIntSortStruct1k(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		d := IntSlice{data}
		b.StartTimer()
		MyIntSortStruct(d)
		b.StopTimer()
	}
}

func BenchmarkMyIntSortPtr1k(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		d := &IntSlicePtr{data}
		b.StartTimer()
		MyIntSortPtr(d)
		b.StopTimer()
	}
}
