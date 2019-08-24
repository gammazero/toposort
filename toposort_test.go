package toposort

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {
	// Test sorting a DAG.
	t.Log("Sorting graph:")
	t.Log("A--> B--> D--> E <---F")
	t.Log("|         ^          |")
	t.Log("|         |          |")
	t.Log("+-------> C <--------+")

	sorted, err := Toposort([]Edge{
		{"B", "D"}, {"D", "E"}, {"A", "B"}, {"A", "C"},
		{"C", "D"}, {"F", "C"}, {"F", "E"}})

	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}

	ss := make([]string, len(sorted), len(sorted))
	for i, n := range sorted {
		ss[i] = n.(string)
	}
	sortedString := strings.Join(ss, "")

	// Check that all values are present in sorted list.
	for _, v := range []string{"A", "B", "C", "D", "E", "F"} {
		if !strings.Contains(sortedString, v) {
			t.Fatal("missing node from sorted result")
		}
	}

	iA := strings.Index(sortedString, "A")
	iB := strings.Index(sortedString, "B")
	iC := strings.Index(sortedString, "C")
	iD := strings.Index(sortedString, "D")
	iE := strings.Index(sortedString, "E")
	iF := strings.Index(sortedString, "F")
	if !((iA < iB) && (iA < iC) && (iB < iD) && (iD < iE) && (iC < iD) && (iF < iC)) {
		t.Fatal("items are not correctly sorted")
	}
	t.Log("Sorted correctly:", sorted)
}

func TestOneSided(t *testing.T) {
	sorted, err := Toposort([]Edge{
		{"A", "B"}, {"A", "C"}, {"A", "D"}, {"A", "E"}, {"A", "F"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(sorted) != 6 {
		t.Fatal("Missing expected nodes in sorted list:", sorted)
	}
}

func TestBadNode(t *testing.T) {
	_, err := Toposort([]Edge{{"X", "X"}})
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestCycle(t *testing.T) {
	// Test sorting a DAG with a cycle.
	t.Log("Sorting graph with cycle:")
	t.Log("          +---------------+")
	t.Log("          |               |")
	t.Log("A--> B--> D--> E <---F <--+")
	t.Log("|         ^          |")
	t.Log("|         |          |")
	t.Log("+-------> C <--------+")
	// There is a cycle: D->F->C->D
	_, err := Toposort([]Edge{
		{"B", "D"}, {"D", "E"}, {"A", "B"}, {"A", "C"},
		{"C", "D"}, {"F", "C"}, {"F", "E"}, {"D", "F"}})
	if err == nil {
		t.Fatal("Toposort failed to detect cycle")
	}
	t.Log(err)
}

// Check that results are correct with many different edge orderings.
func TestBumstead(t *testing.T) {
	// Edges are (x, y) where x depends on y.  In a DAG: y-->x
	clothing := []Edge{
		{"jacket", "tie"}, {"jacket", "belt"},
		{"tie", "shirt"},
		{"belt", "shirt"}, {"belt", "pants"},
		{"pants", "undershorts"},
		{"shoes", "pants"}, {"shoes", "undershorts"}, {"shoes", "socks"},
		{"watch", nil}}

	t.Log("Sorting Professor Bumstead's cloths:")
	sorted, err := ToposortR(clothing)
	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}

	err = validateClothing(clothing, sorted)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Sorted correctly:", sorted)

	rand.Seed(time.Now().Unix())
	for i := 0; i < 37; i++ {
		shuffle(clothing)
		sorted, err := ToposortR(clothing)
		if err != nil {
			t.Fatal("Toposort returned error:", err)
		}
		err = validateClothing(clothing, sorted)
		if err != nil {
			t.Fatal(err)
		}
		t.Log("Sorted correctly:", sorted)
	}
}

// Check that cycle is always detected, with different edge orderings.
func TestBumsteadCycle(t *testing.T) {
	// Edges are (x, y) where x depends on y.  In a DAG: y-->x
	clothing := []Edge{
		{"jacket", "tie"}, {"jacket", "belt"},
		{"tie", "shirt"},
		{"undershorts", "shoes"},
		{"belt", "shirt"}, {"belt", "pants"},
		{"pants", "undershorts"},
		{"shoes", "pants"}, {"shoes", "undershorts"}, {"shoes", "socks"},
		{"watch", nil}}

	for i := 0; i < 7; i++ {
		_, err := ToposortR(clothing)
		if err == nil {
			t.Fatal("failed to detect cycle")
		}
		shuffle(clothing)
	}
}

// Test with multiple edges between the same vertexes.
func TestMultiLink(t *testing.T) {
	t.Log("Sorting multilink")
	sorted, err := Toposort([]Edge{
		{"A", "B"}, {"A", "B"}, {"A", "B"}, {"A", "C"},
		{"A", "C"}, {"B", "C"}, {"B", "C"}, {"B", "C"},
		{"C", "D"}, {"C", "D"}, {"B", "D"}, {"B", "D"}})
	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}
	if len(sorted) != 4 {
		t.Fatal("missing vertexes in output, have", sorted)
	}
	if sorted[0] != "A" || sorted[1] != "B" || sorted[2] != "C" {
		t.Fatal("wrong order")
	}
	t.Log("Multilink sorted:", sorted)
}

// Large
func TestLarge(t *testing.T) {
	graph := []Edge{
		{"A", "B"}, {"A", "C"}, {"C", "B"}, {"C", "E"},
		{"B", "E"}, {"B", "D"}, {"B", "G"}, {"E", "K"},
		{"E", "D"}, {"C", "D"}, {"D", "K"}, {"D", "F"},
		{"G", "F"}, {"F", "K"}, {"F", "J"}, {"F", "I"},
		{"F", "H"}, {"K", "L"}, {"L", "M"}, {"M", "J"},
		{"I", "N"}, {"J", "N"}, {"O", "K"}, {"K", "P"},
		{"Q", "R"}, {"R", "C"}, {"R", "S"}, {"S", "C"},
		{"T", "U"}, {"U", "C"}, {"U", "V"}, {"V", "W"},
		{"W", "Q"}, {"X", "C"}, {nil, "Y"}, {"Z", nil}}

	t.Log("Sorting graph:")
	rand.Seed(time.Now().Unix())
	for i := 0; i < 37; i++ {
		shuffle(graph)
		sorted, err := Toposort(graph)
		if err != nil {
			t.Fatal("Toposort returned error:", err)
		}
		if len(sorted) != 26 {
			t.Fatal("missing vertexes, have", sorted)
		}
		t.Log("Sorted correctly:", sorted)
	}
}

func TestSingleNodeEdge(t *testing.T) {
	sorted, err := Toposort([]Edge{{"A", "B"}, {"A", "C"}, {"A", nil}})
	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}
	if len(sorted) != 3 {
		t.Fatal("missing vertexes, have", sorted)
	}
	if sorted[0] != "A" {
		t.Fatal("wrong order")
	}
}

func BenchmarkToposort(b *testing.B) {
	graph := []Edge{
		{"A", "B"}, {"A", "C"}, {"C", "B"}, {"C", "E"},
		{"B", "E"}, {"B", "D"}, {"B", "G"}, {"E", "K"},
		{"E", "D"}, {"C", "D"}, {"D", "K"}, {"D", "F"},
		{"G", "F"}, {"F", "K"}, {"F", "J"}, {"F", "I"},
		{"F", "H"}, {"K", "L"}, {"L", "M"}, {"M", "J"},
		{"I", "N"}, {"J", "N"}, {"Z", "A"}, {"Y", "A"},
		{"Y", "Z"}, {"W", "A"}, {"W", "Y"}}

	rand.Seed(time.Now().Unix())
	shuffle(graph)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Toposort(graph)
	}
}

// validateClothing checks that Professor Bumstead dressed himself properly.
func validateClothing(clothing []Edge, sorted []interface{}) error {
	// Make sure he is wearing all his clothing items.
	for _, c := range clothing {
		child, parent := c[0], c[1]
		// Check the child is in sorted
		var found bool
		for i := 0; i < len(sorted); i++ {
			if child == sorted[i] {
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprint("missing item: ", child))
		}

		// Check that parent is nil or in sorted.
		found = false
		for i := 0; i < len(sorted); i++ {
			if parent == nil || parent == sorted[i] {
				found = true
				break
			}
		}
		if !found {
			return errors.New(fmt.Sprint("missing item: ", parent))
		}
	}

	// Make sure he put his clothes on in the right order.
	iUndershorts := index(sorted, "undershorts")
	iPants := index(sorted, "pants")
	iBelt := index(sorted, "belt")
	iJacket := index(sorted, "jacket")
	iShirt := index(sorted, "shirt")
	iTie := index(sorted, "tie")
	iSocks := index(sorted, "socks")
	iShoes := index(sorted, "shoes")
	//iWatch := sorted.index("watch")
	if !((iUndershorts < iPants) &&
		(iUndershorts < iShoes) &&
		(iPants < iShoes) &&
		(iPants < iBelt) &&
		(iShirt < iBelt) &&
		(iShirt < iTie) &&
		(iTie < iJacket) &&
		(iSocks < iShoes)) {
		errors.New("clothing items are not correctly sorted")
	}
	return nil
}

func index(slice []interface{}, value string) int {
	for p, v := range slice {
		if v != nil && v.(string) == value {
			return p
		}
	}
	return -1
}

func shuffle(x []Edge) {
	for i := len(x) - 1; i >= 0; i-- {
		// Pick an element in x[:i+1] with which to exchange x[i]
		j := int(rand.Float32() * float32(i+1))
		x[i], x[j] = x[j], x[i]
	}
}
