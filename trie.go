package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type ASNDataTrie struct {
	Root *N
}

type N struct {
	One  *N
	Zero *N
	ASN  int
}

func parseMask(row string) (prefixes []*net.IPNet, asn int) {
	data := strings.Split(row, "\t")

	asn, err := strconv.Atoi(data[2])
	if err != nil {
		log.Fatal(err)
	}

	return ipv4RangeToCidr(data[0], data[1]), asn
}

func (data *ASNDataTrie) Insert(prefix *net.IPNet, asn int) {
	n := data.Root
	ip := prefix.IP
	mask, _ := prefix.Mask.Size()
	for i := 0; i < mask; i++ {
		n = n.getOrCreateChild(bitAt(&ip, i))
	}
	n.ASN = asn
}

func bitAt(ip *net.IP, i int) byte {
	octet := (*ip)[i/8]
	bit := (octet >> (7 - i%8)) & 1
	return bit
}

func (data *ASNDataTrie) GetASNForIp(ipString string) int {
	ip := net.ParseIP(ipString).To4()
	n := data.Root

	for i := 0; i < len(ip)*8; i++ {
		bit := bitAt(&ip, i)

		if bit == 1 && n.One != nil {
			n = n.One
		} else if bit == 0 && n.Zero != nil {
			n = n.Zero
		} else {
			return n.ASN
		}
	}

	return n.ASN
}

func (n *N) getOrCreateChild(bit byte) *N {
	if bit == 1 {
		if n.One == nil {
			n.One = new(N)
		}
		return n.One
	} else {
		if n.Zero == nil {
			n.Zero = new(N)
		}
		return n.Zero
	}
}

// --- loading

func fromFile(fileName string) ASNDataTrie {
	asnData, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal("Could not read file!")
	}
	d := ASNDataTrie{Root: new(N)}
	for _, row := range strings.Split(string(asnData), "\n") {
		if row == "" {
			continue
		}
		prefixes, asn := parseMask(row)
		for _, p := range prefixes {
			d.Insert(p, asn)
		}
	}

	return d
}

func (data *ASNDataTrie) marshall(fileName string) {
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
