package mysort

type IntSlice struct {
	data []int
}

func (is IntSlice) Len() int {
	return len(is.data)
}

func (is IntSlice) Less(i, j int) bool {
	return is.data[i] < is.data[j]
}

func (is IntSlice) Swap(i, j int) {
	is.data[i], is.data[j] = is.data[j], is.data[i]
}

func insertionSortStruct(data IntSlice, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

func siftDownStruct(data IntSlice, lo, hi, first int) {
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

func heapSortStruct(data IntSlice, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownStruct(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDownStruct(data, lo, i, first)
	}
}

func medianOfThreeStruct(data IntSlice, a, b, c int) {
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

func swapRangeStruct(data IntSlice, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}

func doPivotStruct(data IntSlice, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeStruct(data, lo, lo+s, lo+2*s)
		medianOfThreeStruct(data, m, m-s, m+s)
		medianOfThreeStruct(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeStruct(data, lo, m, hi-1)
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
	swapRangeStruct(data, lo, b-n, n)
	n = min(hi-d, d-c)
	swapRangeStruct(data, c, hi-n, n)
	return lo + b - a, hi - (d - c)
}

func quickSortStruct(data IntSlice, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortStruct(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotStruct(data, a, b)
		if mlo-a < b-mhi {
			quickSortStruct(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSortStruct(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortStruct(data, a, b)
	}
}

func MyIntSortStruct(data IntSlice) {
	n := data.Len()
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortStruct(data, 0, n, maxDepth)
}
