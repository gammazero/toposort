package toposort

import (
	"math/rand"
	"slices"
	"testing"
)

func TestSimple(t *testing.T) {
	// Test sorting a DAG.
	t.Log("Sorting graph:")
	t.Log("A--> B--> D--> E <---F")
	t.Log("|         ^          |")
	t.Log("|         |          |")
	t.Log("+-------> C <--------+")

	sorted, err := Toposort([]Edge[string]{
		{"B", "D"}, {"D", "E"}, {"A", "B"}, {"A", "C"},
		{"C", "D"}, {"F", "C"}, {"F", "E"}})

	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}

	// Check that all values are present in sorted list.
	for _, v := range []string{"A", "B", "C", "D", "E", "F"} {
		if !slices.Contains(sorted, v) {
			t.Fatal("missing node from sorted result")
		}
	}

	iA := slices.Index(sorted, "A")
	iB := slices.Index(sorted, "B")
	iC := slices.Index(sorted, "C")
	iD := slices.Index(sorted, "D")
	iE := slices.Index(sorted, "E")
	iF := slices.Index(sorted, "F")
	if (iA >= iB) || (iA >= iC) || (iB >= iD) || (iD >= iE) || (iC >= iD) || (iF >= iC) {
		t.Fatal("items are not correctly sorted")
	}
	t.Log("Sorted correctly:", sorted)
}

func TestOneSided(t *testing.T) {
	sorted, err := Toposort([]Edge[string]{
		{"A", "B"}, {"A", "C"}, {"A", "D"}, {"A", "E"}, {"A", "F"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(sorted) != 6 {
		t.Fatal("Missing expected nodes in sorted list:", sorted)
	}
}

func TestBadNode(t *testing.T) {
	_, err := Toposort([]Edge[string]{{"X", "X"}})
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
	_, err := Toposort([]Edge[string]{
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
	clothing := []Edge[string]{
		{"jacket", "tie"}, {"jacket", "belt"},
		{"tie", "shirt"},
		{"belt", "shirt"}, {"belt", "pants"},
		{"pants", "undershorts"},
		{"shoes", "pants"}, {"shoes", "undershorts"}, {"shoes", "socks"},
		{"watch", ""}}

	t.Log("Sorting Professor Bumstead's cloths:")
	sorted, err := ToposortR(clothing)
	if err != nil {
		t.Fatal("Toposort returned error:", err)
	}

	validateClothing(t, clothing, sorted)
	t.Log("Sorted correctly:", sorted)

	for range 37 {
		shuffle(clothing)
		sorted, err := ToposortR(clothing)
		if err != nil {
			t.Fatal("Toposort returned error:", err)
		}
		validateClothing(t, clothing, sorted)
		t.Log("Sorted correctly:", sorted)
	}
}

// Check that cycle is always detected, with different edge orderings.
func TestBumsteadCycle(t *testing.T) {
	// Edges are (x, y) where x depends on y.  In a DAG: y-->x
	clothing := []Edge[string]{
		{"jacket", "tie"}, {"jacket", "belt"},
		{"tie", "shirt"},
		{"undershorts", "shoes"},
		{"belt", "shirt"}, {"belt", "pants"},
		{"pants", "undershorts"},
		{"shoes", "pants"}, {"shoes", "undershorts"}, {"shoes", "socks"},
		{"watch", ""}}

	for range 7 {
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
	sorted, err := Toposort([]Edge[string]{
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
	graph := []Edge[string]{
		{"A", "B"}, {"A", "C"}, {"C", "B"}, {"C", "E"},
		{"B", "E"}, {"B", "D"}, {"B", "G"}, {"E", "K"},
		{"E", "D"}, {"C", "D"}, {"D", "K"}, {"D", "F"},
		{"G", "F"}, {"F", "K"}, {"F", "J"}, {"F", "I"},
		{"F", "H"}, {"K", "L"}, {"L", "M"}, {"M", "J"},
		{"I", "N"}, {"J", "N"}, {"O", "K"}, {"K", "P"},
		{"Q", "R"}, {"R", "C"}, {"R", "S"}, {"S", "C"},
		{"T", "U"}, {"U", "C"}, {"U", "V"}, {"V", "W"},
		{"W", "Q"}, {"X", "C"}, {"", "Y"}, {"Z", ""}}

	t.Log("Sorting graph:")
	for range 37 {
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
	sorted, err := Toposort([]Edge[string]{{"A", "B"}, {"A", "C"}, {"A", ""}})
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
	graph := []Edge[string]{
		{"A", "B"}, {"A", "C"}, {"C", "B"}, {"C", "E"},
		{"B", "E"}, {"B", "D"}, {"B", "G"}, {"E", "K"},
		{"E", "D"}, {"C", "D"}, {"D", "K"}, {"D", "F"},
		{"G", "F"}, {"F", "K"}, {"F", "J"}, {"F", "I"},
		{"F", "H"}, {"K", "L"}, {"L", "M"}, {"M", "J"},
		{"I", "N"}, {"J", "N"}, {"Z", "A"}, {"Y", "A"},
		{"Y", "Z"}, {"W", "A"}, {"W", "Y"}}

	shuffle(graph)
	b.ResetTimer()
	for range b.N {
		_, _ = Toposort(graph)
	}
}

// validateClothing checks that Professor Bumstead dressed himself properly.
func validateClothing(t *testing.T, clothing []Edge[string], sorted []string) {
	t.Helper()
	t.Log("Sorted:", sorted)
	// Make sure he is wearing all his clothing items.
	for _, c := range clothing {
		child, parent := c[0], c[1]
		// Check the child is in sorted
		if !slices.Contains(sorted, child) {
			t.Fatalf("missing child item: %q", child)
		}

		// Check that parent is "" or in sorted.
		if !slices.ContainsFunc(sorted, func(s string) bool {
			return parent == "" || s == parent
		}) {
			t.Fatalf("missing prent item: %q", parent)
		}
	}

	// Make sure he put his clothes on in the right order.
	iUndershorts := slices.Index(sorted, "undershorts")
	iPants := slices.Index(sorted, "pants")
	iBelt := slices.Index(sorted, "belt")
	iJacket := slices.Index(sorted, "jacket")
	iShirt := slices.Index(sorted, "shirt")
	iTie := slices.Index(sorted, "tie")
	iSocks := slices.Index(sorted, "socks")
	iShoes := slices.Index(sorted, "shoes")
	//iWatch := sorted.index("watch")
	if (iUndershorts >= iPants) ||
		(iUndershorts >= iShoes) ||
		(iPants >= iShoes) ||
		(iPants >= iBelt) ||
		(iShirt >= iBelt) ||
		(iShirt >= iTie) ||
		(iTie >= iJacket) ||
		(iSocks >= iShoes) {
		t.Fatal("clothing items are not correctly sorted")
	}
}

func index(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func shuffle(x []Edge[string]) {
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
}
