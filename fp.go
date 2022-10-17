package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type ASNRecord struct {
	StartIp uint32
	ASN     int
}

type ASNData struct {
	Records []ASNRecord
}

func parseIp(ip string) uint32 {
	ipBytes := [4]byte{}
	ipParts := strings.Split(ip, ".")
	for i, part := range ipParts {
		s, err := strconv.ParseInt(part, 10, 32)
		if err != nil {
			log.Fatalf("Could not parse int from %s\n", part)
		}
		ipBytes[i] = byte(s)
	}
	return binary.BigEndian.Uint32(ipBytes[:])
}

// accepts "start end asn"
func parseRecord(row string) ASNRecord {
	parts := strings.Split(row, "\t")
	startIp := parseIp(parts[0])
	asn, err := strconv.Atoi(parts[2])
	if err != nil {
		log.Fatalf("Could not parse ASN from %s\n", parts[2])
	}
	return ASNRecord{StartIp: startIp, ASN: asn}
}

func buildSearchAlgorithm(data string) ASNData {
	var parsedRecords []ASNRecord
	for _, row := range strings.Split(data, "\n") {
		if row == "" { // empty line
			continue
		}
		parsedRecords = append(parsedRecords, parseRecord(row))
	}
	sort.Slice(parsedRecords, func(i, j int) bool {
		return parsedRecords[i].StartIp < parsedRecords[j].StartIp
	})
	return ASNData{parsedRecords}
}

func (data *ASNData) getAsn(ip string) int {
	dataLen := len(data.Records)
	s := parseIp(ip)
	i, j := 0, dataLen
	for i < j {
		h := int(uint(i+j) >> 1)
		a := data.Records[h]

		if a.StartIp <= s {
			if h+1 == dataLen || data.Records[h+1].StartIp > s {
				return a.ASN
			}
			i = h + 1
		} else {
			j = h
		}
	}
	return data.Records[i].ASN
}

// --- data loading and marshalling

func loadFromDataFile(fileName string) ASNData {
	asnData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not read file!")
	}
	return buildSearchAlgorithm(string(asnData))
}

func loadSaved(fileName string) ASNData {
	saved, err := os.ReadFile(fileName)
	dec := gob.NewDecoder(bytes.NewBuffer(saved))
	if err != nil {
		log.Fatal(err)
	}
	var data ASNData
	err = dec.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func (data *ASNData) marshall(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Couldn't open file")
	}
	defer f.Close()

	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err = enc.Encode(data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}
