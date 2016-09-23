package mysort

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func insertionSort(data []int, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && (data[j] < data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func siftDown(data []int, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && (data[first+child] < data[first+child+1]) {
			child++
		}
		if !(data[first+root] < data[first+child]) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func heapSort(data []int, a, b int) {
	first := a
	lo := 0
	hi := b - a
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDown(data, lo, i, first)
	}
}

func medianOfThree(data []int, a, b, c int) {
	m0 := b
	m1 := a
	m2 := c
	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
	if data[m2] < data[m1] {
		data[m2], data[m1] = data[m1], data[m2]
	}
	if data[m1] < data[m0] {
		data[m1], data[m0] = data[m0], data[m1]
	}
}

func swapRange(data []int, a, b, n int) {
	for i := 0; i < n; i++ {
		data[a+i], data[b+i] = data[b+i], data[a+i]
	}
}

func doPivot(data []int, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	medianOfThree(data, lo, m, hi-1)
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data[b] < data[pivot] { // data[b] < pivot
				b++
			} else if !(data[pivot] < data[b]) { // data[b] = pivot
				data[a], data[b] = data[b], data[a]
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data[pivot] < data[c-1] { // data[c-1] > pivot
				c--
			} else if !(data[c-1] < data[pivot]) { // data[c-1] = pivot
				data[c-1], data[d-1] = data[d-1], data[c-1]
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}
	n := min(b-a, a-lo)
	swapRange(data, lo, b-n, n)
	n = min(hi-d, d-c)
	swapRange(data, c, hi-n, n)
	return lo + b - a, hi - (d - c)
}

func quickSort(data []int, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(data, a, b)
	}
}

func MyIntSort(data []int) {
	n := len(data)
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(data, 0, n, maxDepth)
}
