# Benchmark for ASN lookup

- download dataset from https://iptoasn.com/
- [fp.go](fp.go) uses sorted list of integers for an ASN lookup
- [trie.go](trie.go) then uses binary [Trie](https://en.wikipedia.org/wiki/Trie) for a lookup
- [intersections.go](intersections.go) is a file that shows that the dataset has overlapping IP ranges