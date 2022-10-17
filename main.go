package main

func main() {
	data := loadFromDataFile("./ip2asn-v4.tsv")
	println(data.getAsn("1.0.0.1"))

	data.marshall("saved.bin")

	data = loadSaved("saved.bin")
	println(data.getAsn("1.0.0.1"))
}
