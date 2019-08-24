package toposort

import (
	"fmt"
	"log"
)

// This example shows dependency-sorting articles of clothing to compute a
// solution for Professor Bumstead to get dressed correctly.  Articles of
// clothing are arranged into pairs where the first item in a pair must be put
// on after the second item in the pair.
func Example_clothes() {
	fmt.Println("Dependency-sorting Professor Bumstead's cloths:")

	// Edges are (x, y) where x depends on y.  In other words, y must be done
	// before x.  In a DAG: y --> x.  So ToposortR is called for this reversed
	// order.
	sorted, err := ToposortR([]Edge{
		{"jacket", "tie"}, {"jacket", "belt"},
		{"tie", "shirt"},
		{"belt", "shirt"}, {"belt", "pants"},
		{"pants", "undershorts"},
		{"shoes", "pants"}, {"shoes", "undershorts"}, {"shoes", "socks"},
		{"watch", nil}})
	if err != nil {
		log.Fatal("Toposort returned error:", err)
	}
	fmt.Println("Sorted correctly:", sorted)
}
