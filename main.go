package main

import (
	"fmt"
	"net"
)

func runFp(ip string) {
	data := loadFromDataFile("./ip2asn-v4.tsv")
	fmt.Printf("%d\n", data.GetASNForIp(ip))
}

func runTrie(ip string) {
	data := fromFile("./ip2asn-v4.tsv")
	fmt.Printf("%d\n", data.GetASNForIp(ip))

	data.marshall("./saved2.bin")
}

func test() {
	//ip := net.ParseIP("1.255.255.0").To4()
	//
	//fmt.Printf("%08b\n", ip)
	//
	//for i := 0; i < len(ip)*8; i++ {
	//	octet := ip[i/8]
	//	bit := (octet >> (7 - i%8)) & 1
	//	print(bit)
	//}

	ip1 := net.ParseIP("1.0.0.0").To4()
	ip2 := net.ParseIP("1.0.0.255").To4()

	mask := make([]byte, len(ip1))
	for i := range mask {
		mask[i] = ip1[i] ^ (^ip2[i])
	}
	ipnet := net.IPNet{IP: ip1, Mask: mask}
	size, bits := ipnet.Mask.Size()

	fmt.Printf("%08b\n", mask)
	fmt.Printf("%08b\n", ip1)
	fmt.Printf("%08b\n", ip2)
	fmt.Printf("%08b\n", ip1.Mask(mask))
	fmt.Printf("%d %d\n", size, bits)

}

func main() {
	//ip := "1.0.213.2"
	//runFp(ip)
	//runTrie(ip)
	//1.0.223.0	1.0.255.255
	inet, _ := parseMask("1.0.223.0\t1.0.255.255\t123")
	fmt.Printf("%08b\n", *inet)
	ones, bits := inet.Mask.Size()
	fmt.Printf("%d %d\n", ones, bits)
}
