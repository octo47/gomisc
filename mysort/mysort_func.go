package mysort

func insertionSortFunc(data []int, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && Less(data, j, j-1); j-- {
			Swap(data, j, j-1)
		}
	}
}

func siftDownFunc(data []int, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && Less(data, first+child, first+child+1) {
			child++
		}
		if !Less(data, first+root, first+child) {
			return
		}
		Swap(data, first+root, first+child)
		root = child
	}
}

func heapSortFunc(data []int, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDownFunc(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		Swap(data, first, first+i)
		siftDownFunc(data, lo, i, first)
	}
}

func medianOfThreeFunc(data []int, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	if Less(data, m1, m0) {
		Swap(data, m1, m0)
	}
	if Less(data, m2, m1) {
		Swap(data, m2, m1)
	}
	if Less(data, m1, m0) {
		Swap(data, m1, m0)
	}
}

func swapRangeFunc(data []int, a, b, n int) {
	for i := 0; i < n; i++ {
		Swap(data, a+i, b+i)
	}
}

func doPivotFunc(data []int, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThreeFunc(data, lo, lo+s, lo+2*s)
		medianOfThreeFunc(data, m, m-s, m+s)
		medianOfThreeFunc(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThreeFunc(data, lo, m, hi-1)
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if Less(data, b, pivot) { // data[b] < pivot
				b++
			} else if !Less(data, pivot, b) { // data[b] = pivot
				Swap(data, a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if Less(data, pivot, c-1) { // data[c-1] > pivot
				c--
			} else if !Less(data, c-1, pivot) { // data[c-1] = pivot
				Swap(data, c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		Swap(data, b, c-1)
		b++
		c--
	}
	n := min(b-a, a-lo)
	swapRangeFunc(data, lo, b-n, n)
	n = min(hi-d, d-c)
	swapRangeFunc(data, c, hi-n, n)
	return lo + b - a, hi - (d - c)
}

func quickSortFunc(data []int, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSortFunc(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivotFunc(data, a, b)
		if mlo-a < b-mhi {
			quickSortFunc(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSortFunc(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSortFunc(data, a, b)
	}
}

func MyIntSortFunc(data []int) {
	n := Len(data)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSortFunc(data, 0, n, maxDepth)
}

func Len(data []int) int {
	return len(data)
}

func Less(data []int, i, j int) bool {
	return data[i] < data[j]
}

func Swap(data []int, i, j int) {
	data[i], data[j] = data[j], data[i]
}
