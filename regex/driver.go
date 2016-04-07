package regex
import (
	"fmt"
	"time"
	"bufio"
	"log"
	"os"
	"runtime/pprof"
)

func RunBenchmark(matcher func(string) bool, name string) {
	f, err := os.Create(fmt.Sprint("%s.prof", name))
	if err != nil {
		log.Fatal(err)
	}
 	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	start := time.Now()
	log.Printf("Started %s: %v", name, start)
	inputFile, err := os.Open("log-sample.txt")
	if err != nil {
		log.Fatal("Error opening input file:", err)
	}
	defer inputFile.Close()
	reader := bufio.NewReaderSize(bufio.NewReader(inputFile), 1 << 20) // 1mb
	for line, err := reader.ReadString('\n'); err==nil;
		line, err = reader.ReadString('\n'){
		if matcher(line) {
			fmt.Printf("Match: %v\n", string(line))
		}
	}
	fmt.Printf("Elapsed time: %v\n", time.Since(start))
}


