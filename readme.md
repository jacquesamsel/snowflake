# snowflake
Snowflake is a fast, goroutine-safe unique ID generator built for distributed systems

[![GoDoc](https://godoc.org/github.com/Dextication/snowflake?status.svg)](https://godoc.org/github.com/Dextication/snowflake)
[![Go report](http://goreportcard.com/badge/dextication/snowflake)](http://goreportcard.com/report/dextication/snowflake)
[![Coverage](http://gocover.io/_badge/github.com/Dextication/snowflake)](https://gocover.io/github.com/Dextication/snowflake)
## Key concepts
### Snowflake
Snowflakes are `int64`s. `uint64` is not used for interoperability reasons, and to allow
for correct sorting. Snowflakes can be naturally sorted by time. The smaller the snowflake, the earlier it was created.
A snowflake can be broken down into these parts:

- 1 unused bit
- timestamp (ms) since custom epoch (customisable size)
- node id (customisable size)
- counter (customisable size)

Due to the space restrictions of an int64, the timestamp, node id and counter must share 63 bits. It is recommended to
plan for the future: the amount of bits allocated for the timestamp should be as large as possible. A simple formula to
determine how many bits to allocate is this: `ceil(log2(years * 3,154e+10))`. `years` represents how long in the future
you wish for the ids to generate properly.
### Node
A `Node` generates IDs - it is recommended to make use of one `Node` per machine, although there are situations where it
may also be acceptable to use more.

A `Node` coordinates the generation of IDs. It maintains a counter to ensure that IDs are unique. A `Node` also stores
a custom epoch, which should be defined by the user.

## Examples
### ID Generation
The following example allocates:
- 43 bits timestamp (valid for 278.7 years + customEpoch)
- 10 bits node id (1024 nodes)
- 10 bits counter (max. 1024 snowflakes/ms/node)
```go 
node, err := snowflake.NewNode(0, customEpoch, 43, 10, 10)
if err != nil {
    panic(err)
}

id := node.Generate()
```
### Base64
Snowflakes can be converted to base64 (mathematical base using `0-9a-zA-Z+/`). This is more space-efficient than base64 because it excludes
padding.
```go 
// Snowflake to base64
s := snowflake.Snowflake(4294967296)
fmt.Println(s.Base64()) // 400000

// Base64 to snowflake
s = snofwlake.ParseBase64("400000")
fmt.Println(s) // 4294967296
```

## Benchmarks
To test the performance on your local machine, clone this repository and run `go test -bench=.`.
```
cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
BenchmarkNode_Generate-12       51937034                23.24 ns/op            0 B/op          0 allocs/op
BenchmarkParseBase64-12         34151838                35.10 ns/op            0 B/op          0 allocs/op
```