# Benchmark for ASN lookup

- download dataset from https://iptoasn.com/
- [fp.go](fp.go) uses sorted list of integers for an ASN lookup
- [trie.go](trie.go) then uses binary [Trie](https://en.wikipedia.org/wiki/Trie) for a lookup
- [intersections.go](intersections.go) is a file that shows that the dataset has overlapping IP ranges

## Observations

The datastructures in the benchmark were not optimized at all, I just wanted to see how much worse the Trie/List are
without tuning them up. Moreover, I didn't expect that there are overlapping IP ranges in the dataset (183), so this
case is not handled as well. I presume supporting this in the both implementations would make Trie to have smaller
stored size than the Sorted List.

I run `main.go` on my M1 Pro mac with `1 000 000` searches for IP

| Method      | Build Time | Search / IP |   Stored Size |
|:------------|:----------:|------------:|--------------:|
| Sorted List |  101 mls   |      145 ns | 5000273 bytes |
| Trie        |  562 mls   |      201 ns | 5907751 bytes |

## Conclusion

Naive implementation of Sorted List is better than naive implementation fo Trie. 
