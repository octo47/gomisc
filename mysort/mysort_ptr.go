package mysort

type IntSlicePtr struct {
	data []int
}

func (is *IntSlicePtr) Len() int {
	return len(is.data)
}

func (is *IntSlicePtr) Less(i, j int) bool {
	return is.data[i] < is.data[j]
}

func (is *IntSlicePtr) Swap(i, j int) {
	is.data[i], is.data[j] = is.data[j], is.data[i]
}

func insertionSortPtr(data *IntSlicePtr, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func siftDownPtr(data *IntSlicePtr, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSortPtr(data *IntSlicePtr, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownPtr(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDownPtr(data, lo, i, first)
	}
}

func medianOfThreePtr(data *IntSlicePtr, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
	if data.Less(m2, m1) {
		data.Swap(m2, m1)
	}
	if data.Less(m1, m0) {
		data.Swap(m1, m0)
	}
}

func swapRangePtr(data *IntSlicePtr, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}

func doPivotPtr(data *IntSlicePtr, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreePtr(data, lo, lo+s, lo+2*s)
		medianOfThreePtr(data, m, m-s, m+s)
		medianOfThreePtr(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreePtr(data, lo, m, hi-1)
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data.Less(b, pivot) { // data[b] < pivot
				b++
			} else if !data.Less(pivot, b) { // data[b] = pivot
				data.Swap(a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data.Less(pivot, c-1) { // data[c-1] > pivot
				c--
			} else if !data.Less(c-1, pivot) { // data[c-1] = pivot
				data.Swap(c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		data.Swap(b, c-1)
		b++
		c--
	}
	n := min(b-a, a-lo)
	swapRangePtr(data, lo, b-n, n)
	n = min(hi-d, d-c)
	swapRangePtr(data, c, hi-n, n)
	return lo + b - a, hi - (d - c)
}

func quickSortPtr(data *IntSlicePtr, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortPtr(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotPtr(data, a, b)
		if mlo-a < b-mhi {
			quickSortPtr(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSortPtr(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortPtr(data, a, b)
	}
}

func MyIntSortPtr(data *IntSlicePtr) {
	n := data.Len()
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortPtr(data, 0, n, maxDepth)
}
