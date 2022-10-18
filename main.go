package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

func runFp(ips []string) {
	start := time.Now()

	data := loadFromDataFile("./ip2asn-v4.tsv")

	log.Printf("SORT - build: %s\n", time.Since(start))

	start = time.Now()

	for _, ip := range ips {
		data.GetASNForIp(ip)
	}
	elapsed := time.Since(start)
	nano := elapsed.Nanoseconds() / int64(len(ips))
	log.Printf("SORT - search: total %s, per IP %d ns\n", elapsed, nano)

	file := "./saved_sort.bin"
	data.marshall(file)
	printFileSize(file, "SORT")
}

func runTrie(ips []string) {
	start := time.Now()

	data := fromFile("./ip2asn-v4.tsv")

	log.Printf("TRIE - build: %s\n", time.Since(start))

	start = time.Now()

	for _, ip := range ips {
		data.GetASNForIp(ip)
	}
	elapsed := time.Since(start)
	nano := elapsed.Nanoseconds() / int64(len(ips))
	log.Printf("TRIE - search: total %s, per IP %d ns\n", elapsed, nano)

	file := "./saved_trie.bin"
	data.marshall(file)
	printFileSize(file, "TRIE")
}

func printFileSize(file string, name string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s - size: %d bytes", name, info.Size())
}

func generateDataSet(size int) (ips []string) {
	for i := 0; i < size; i++ {
		ips = append(ips, fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256)))
	}
	return ips
}

func main() {
	ips := generateDataSet(100000)
	runFp(ips)
	runTrie(ips)
}
