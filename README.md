# toposort

[![GoDoc](https://godoc.org/github.com/gammazero/toposort?status.svg)](https://godoc.org/github.com/gammazero/toposort)
[![Build Status](https://github.com/gammazero/toposort/actions/workflows/go.yml/badge.svg)](https://github.com/gammazero/toposort/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gammazero/toposort)](https://goreportcard.com/report/github.com/gammazero/toposort)
[![codecov](https://codecov.io/gh/gammazero/toposort/branch/master/graph/badge.svg)](https://codecov.io/gh/gammazero/toposort)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Topologically sort a directed acyclic graph (DAG) with cycle detection.

This topological sort can be used to put items in dependency order, where each
edge (u, v) has a vertex u that is depended on by v, such that u must come
before v.

The result of completing a topological sort is an ordered slice of vertexes,
where the position of every vertex in the list is before any of its destination
vertexes.

After sorting this graph:
```
 A--> B--> D--> E <---F
 |         ^          |
 |         |          |
 +-------> C <--------+
```
The following conditions must be true concerning the relative position of the
nodes in the returned list of nodes: A<B, A<C, B<D, D<E, C<D, F<C, F<E.  The
slice `[F A C B D E]` is a correct result.

This implementation uses [Kahn's algorithm](https://en.wikipedia.org/wiki/Topological_sorting#Kahn.27s_algorithm).

## Installation

```
$ go get github.com/gammazero/toposort
```

## Example

```go
	sorted, err := Toposort([]Edge{
		{"B", "D"}, {"D", "E"}, {"A", "B"}, {"A", "C"},
		{"C", "D"}, {"F", "C"}, {"F", "E"}})
	if err != nil {
		log.Fatal("Toposort returned error:", err)
	}
	fmt.Println("Sorted correctly:", sorted)
```
