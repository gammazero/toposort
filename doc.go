/*
Topologically sort a directed acyclic graph (DAG) with cycle detection.

A topological sort of a DAG G = (V, E) is a linear ordering of all its vertexes
such that if G contains an edge (u, v), then u appears before v in the
ordering.  In other words u < v.

This topological sort can be used to put items in dependency order, where each
edge (u, v) has a vertex u that is depended on by v, such that u must come
before v.

This sort finds cycles in the graph.  If the graph is determined to have a
cycle, then an error is returned.

The result of completing a topological sort is an ordered slice of vertexes,
where the position of every vertex in the list is before any of its destination
vertexes.

After sorting this graph:

 A--> B--> D--> E <---F
 |         ^          |
 |         |          |
 +-------> C <--------+

The following conditions must be true concerning the relative position of the
nodes in the returned list of nodes: A<B, A<C, B<D, D<E, C<D, F<C, F<E.  The
slice [F A C B D E] is a correct result.

When sorting this graph:

           +---------------+
           |               |
 A--> B--> D--> E <---F <--+
 |         ^          |
 |         |          |
 +-------> C <--------+

Toposort will return with an error stating that a cycle was detected.

Reversing the order of the returned vertices is the same as reversing the
direction of each edge.

This implementation uses Kahn's algorithm.
https://en.wikipedia.org/wiki/Topological_sorting#Kahn.27s_algorithm
*/
package toposort
