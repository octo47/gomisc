$ go test -bench .
PASS
BenchmarkSortInt1K-8	   10000	    150809 ns/op
BenchmarkMyIntSort1K-8	   50000	     34777 ns/op
BenchmarkMyIntSortFunc1K-8	   50000	     51170 ns/op
BenchmarkMyIntSortStruct1k-8	   10000	    143204 ns/op
BenchmarkMyIntSortPtr1k-8	   50000	     56585 ns/op
ok  	_/home/jmarsh/sort_test/mysort	13.993s
