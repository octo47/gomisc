package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func Benchmark_Lock(b *testing.B) {
	var lock sync.RWMutex
	var value int32

	for i := 0; i < b.N; i++ {
		lock.Lock()
		value = 5
		lock.Unlock()
	}
	_ = value
}

func Benchmark_Atomic(b *testing.B) {
	var value int32

	for i := 0; i < b.N; i++ {
		atomic.StoreInt32(&value, 5)
	}
	_ = value
}

func Benchmark_Lock_Goroutine(b *testing.B) {
	var lock sync.RWMutex
	var value int32
	value = 5
	wg := new(sync.WaitGroup)

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			lock.RLock()
			defer lock.RUnlock()
			_ = value
		}(wg)
	}
	wg.Wait()
}

func Benchmark_Atomic_Goroutine(b *testing.B) {
	var value int32
	atomic.StoreInt32(&value, 5)

	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			atomic.LoadInt32(&value)
		}(wg)
	}
	wg.Wait()
}
