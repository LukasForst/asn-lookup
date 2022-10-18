package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Interval struct {
	start uint32
	end   uint32

	originalRange string
	asn           int
}

func findOverlaps(fileName string) int {
	asnData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not read file!")
	}

	var intervals []*Interval
	for _, row := range strings.Split(string(asnData), "\n") {
		if row == "" {
			continue
		}

		parts := strings.Split(row, "\t")
		startIp := parseIp(parts[0])
		endIp := parseIp(parts[1])
		asn, err := strconv.Atoi(parts[2])
		if err != nil {
			log.Fatal(err)
		}
		intervals = append(intervals, &Interval{start: startIp, end: endIp, asn: asn, originalRange: fmt.Sprintf("%s - %s", parts[0], parts[1])})
	}
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].start < intervals[j].start
	})

	overlaps := 0
	for i := 1; i < len(intervals); i++ {
		prev := intervals[i-1]
		curr := intervals[i]

		if prev.asn == curr.asn {
			continue
		}

		if prev.end > curr.start {
			log.Printf("Interval Overlap!\nASN %d has %s\nASN %d has %s\n", prev.asn, prev.originalRange, curr.asn, curr.originalRange)
			overlaps++
		}
	}
	return overlaps
}
