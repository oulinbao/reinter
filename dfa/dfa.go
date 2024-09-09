// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package dfa provides a way to construct deterministic finite automata from non-deterministic finite automata.
package dfa

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/oulinbao/regexinter/nfa"
	"github.com/oulinbao/regexinter/runerange"
)

type Node struct {
	State       int  // state
	Final       bool // final?
	Transitions []T  // transitions

	label    string
	closures []*nfa.Node
}

type T struct {
	RuneRanges []rune // rune ranges
	Node       *Node  // node
}

type context struct {
	state        int
	nodesByLabel map[string]*Node
	closureCache map[*nfa.Node][]*nfa.Node
}

var visited = make(map[*Node]bool)

func (n Node) Print() {
	fmt.Println(fmt.Sprintf("Lable:%s, State: %d, Final: %v, Trans: %v", n.label, n.State, n.Final, n.Transitions))

	for _, t := range n.Transitions {
		if _, exist := visited[t.Node]; exist {
			continue
		}

		visited[t.Node] = true
		fmt.Println(fmt.Sprintf("Trans: %v,rune ranges %v", t, t.RuneRanges))

		t.Node.Print()
	}
}

func (n Node) NextState(r []rune) *Node {
	for _, t := range n.Transitions {
		if runerange.Contains(t.RuneRanges, r) {
			return t.Node
		}
	}

	return nil
}

func NewFromNFA(nfanode *nfa.Node) (*Node, error) {
	ctx := &context{
		nodesByLabel: make(map[string]*Node),
		closureCache: make(map[*nfa.Node][]*nfa.Node),
	}
	node := firstNode(nfanode, ctx)
	constructSubset(node, ctx)
	return node, nil
}

func recursiveClosure(node *nfa.Node, visited map[*nfa.Node]bool) []*nfa.Node {
	if visited[node] {
		return nil
	}
	visited[node] = true

	cls := []*nfa.Node{node}
	for _, t := range node.T {
		if t.R == nil { // ε-transition
			cls = append(cls, recursiveClosure(t.N, visited)...)
		}
	}

	return cls
}

func closure(node *nfa.Node, cache map[*nfa.Node][]*nfa.Node) []*nfa.Node {
	if cache != nil {
		if cls, ok := cache[node]; ok {
			return cls
		}
	}

	visited := make(map[*nfa.Node]bool)
	cls := recursiveClosure(node, visited)

	if cache != nil {
		cache[node] = cls
	}

	return cls
}

func intsToStrings(a []int) []string {
	s := make([]string, 0, len(a))
	for _, i := range a {
		s = append(s, strconv.Itoa(i))
	}
	return s
}

func makeLabel(a []int) string {
	return strings.Join(intsToStrings(a), ",")
}

func labelFromClosure(cls []*nfa.Node) string {
	m := make(map[int]struct{})
	for _, n := range cls {
		m[n.S] = struct{}{}
	}

	states := make([]int, 0, len(m))
	for n := range m {
		states = append(states, n)
	}

	sort.Ints(states)

	return makeLabel(states)
}

func isFinal(cls []*nfa.Node) bool {
	visited := make(map[*nfa.Node]bool)
	for _, n := range cls {
		if canReachFinal(n, visited) {
			return true
		}
	}
	return false
}

func canReachFinal(node *nfa.Node, visited map[*nfa.Node]bool) bool {
	if node.F {
		return true
	}
	if visited[node] {
		return false
	}
	visited[node] = true
	for _, t := range node.T {
		if t.R == nil { // 空转换
			if canReachFinal(t.N, visited) {
				return true
			}
		}
	}
	return false
}

func union(cls ...[]*nfa.Node) []*nfa.Node {
	if len(cls) == 1 {
		return cls[0]
	}

	size := 0
	for _, c := range cls {
		size += len(c)
	}

	m := make(map[*nfa.Node]struct{}, size)
	for _, c := range cls {
		for _, n := range c {
			m[n] = struct{}{}
		}
	}

	a := make([]*nfa.Node, 0, len(m))
	for n := range m {
		a = append(a, n)
	}

	return a
}

func closuresForRange(n *Node, rr []rune, ctx *context) (closures [][]*nfa.Node) {
	for i := range n.closures {
		c := n.closures[i]
		for _, t := range c.T {
			if runerange.Contains(t.R, rr) {
				cls := closure(t.N, ctx.closureCache)
				closures = append(closures, cls)
			}
		}
	}
	return
}

func constructSubset(root *Node, ctx *context) {
	var ranges [][]rune
	for i := range root.closures {
		n := root.closures[i]
		for _, t := range n.T {
			ranges = append(ranges, t.R)
		}
	}
	pairs := runerange.Split(ranges)

	m := make(map[*Node][]rune)

	for i := 0; i < len(pairs); i += 2 {
		cls := union(closuresForRange(root, pairs[i:i+2], ctx)...)
		label := labelFromClosure(cls)
		isFinalState := isFinal(cls)
		var node *Node
		if n, ok := ctx.nodesByLabel[label]; ok {
			node = n
		} else {
			ctx.state++
			node = &Node{
				State:    ctx.state,
				Final:    isFinalState,
				label:    label,
				closures: cls,
			}
			ctx.nodesByLabel[label] = node
			constructSubset(node, ctx)
		}
		m[node] = runerange.Sum(m[node], pairs[i:i+2])
	}

	type transition struct {
		node  *Node
		runes []rune
	}
	var transitions []transition
	for n, rr := range m {
		transitions = append(transitions, transition{n, rr})
	}
	if len(m) > 0 {
		sort.Slice(transitions, func(i, j int) bool {
			if len(transitions[i].node.closures) > 0 && len(transitions[j].node.closures) > 0 {
				return transitions[i].node.closures[0].S < transitions[j].node.closures[0].S
			} else {
				return transitions[i].node.State < transitions[j].node.State
			}

		})
		for _, t := range transitions {
			root.Transitions = append(root.Transitions, T{t.runes, t.node})
		}
	}

}

func firstNode(nfanode *nfa.Node, ctx *context) *Node {
	cls := closure(nfanode, ctx.closureCache)
	label := labelFromClosure(cls)

	ctx.state++
	node := &Node{
		State:    ctx.state,
		Final:    isFinal(cls),
		label:    label,
		closures: cls,
	}
	ctx.nodesByLabel[label] = node

	return node
}
