package toposort

import (
	"fmt"
	"log"
)

// This example shows sorting of a directed acyclic graph.
func Example_graph() {
	// Test sorting a DAG.
	fmt.Println("Sorting graph:")
	fmt.Println("A--> B--> D--> E <---F")
	fmt.Println("|         ^          |")
	fmt.Println("|         |          |")
	fmt.Println("+-------> C <--------+")

	sorted, err := Toposort([]Edge{
		{"B", "D"}, {"D", "E"}, {"A", "B"}, {"A", "C"},
		{"C", "D"}, {"F", "C"}, {"F", "E"}})
	if err != nil {
		log.Fatal("Toposort returned error:", err)
	}
	fmt.Println("Sorted correctly:", sorted)
}
