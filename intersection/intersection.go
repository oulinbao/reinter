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

// This package checks that two DFA has intersection or not by DFS algorithm.
package intersection

import (
	"fmt"
	"log"
	"regexp"

	"github.com/oulinbao/regexinter/dfa"
	"github.com/oulinbao/regexinter/nfa"
	"github.com/oulinbao/regexinter/runerange"
	"github.com/pkg/errors"
)

type CombineNode struct {
	Name        string // state1_state2
	Final       bool
	Node1       *dfa.Node
	Node2       *dfa.Node
	Transitions []T
}

type T struct {
	RuneRanges []rune       // rune ranges
	Node       *CombineNode // node
}

func HasIntersection(expr1, expr2 string) (bool, error) {
	node2, err2 := convert2Dfa(expr2)
	if nil != err2 {
		return false, err2
	}

	node1, err1 := convert2Dfa(expr1)
	if nil != err1 {
		return false, err1
	}

	var nodeMap = make(map[string]*CombineNode)
	firstNode := createNode(node1, node2)
	return dfs(firstNode, nodeMap), nil
}

func convert2Dfa(expr string) (*dfa.Node, error) {
	_, err := regexp.Compile(expr)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("invalid regexp: %s", expr))
	}

	nfaNode, err := nfa.New(expr)
	if err != nil {
		log.Fatal(err)
	}
	return dfa.NewFromNFA(nfaNode)
}

func createNode(node1, node2 *dfa.Node) *CombineNode {
	return &CombineNode{
		Name:  nodeName(node1, node2),
		Node1: node1,
		Node2: node2,
		Final: node1.Final && node2.Final,
	}
}

func dfs(startNode *CombineNode, nodeMap map[string]*CombineNode) bool {
	stack := []*CombineNode{startNode}
	visited := make(map[string]bool)

	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		name := nodeName(node.Node1, node.Node2)

		if visited[name] {
			continue
		}
		visited[name] = true
		nodeMap[name] = node

		if node.Final {
			return true
		}

		ranges := findOverlapRanges(node.Node1.Transitions, node.Node2.Transitions)

		for _, r := range ranges {
			nextNodes1 := findNextNodes(node.Node1, r.Range)
			nextNodes2 := findNextNodes(node.Node2, r.Range)

			for _, next1 := range nextNodes1 {
				for _, next2 := range nextNodes2 {
					tmpNode := createNode(next1, next2)
					tmpName := nodeName(next1, next2)
					if !visited[tmpName] {
						stack = append(stack, tmpNode)
					}
				}
			}
		}
	}

	return false
}

func findNextNodes(node *dfa.Node, r []rune) []*dfa.Node {
	var nextNodes []*dfa.Node
	for _, t := range node.Transitions {
		if runerange.Overlaps(t.RuneRanges, r) {
			nextNodes = append(nextNodes, t.Node)
		}
	}
	return nextNodes
}

type RangeInfo struct {
	Range []rune
	Node1 *dfa.Node
	Node2 *dfa.Node
}

func findOverlapRanges(trans1, trans2 []dfa.T) []RangeInfo {
	var result []RangeInfo
	for _, t1 := range trans1 {
		for _, t2 := range trans2 {
			if runerange.Overlaps(t1.RuneRanges, t2.RuneRanges) {
				overlap := findOverlap(t1.RuneRanges, t2.RuneRanges)
				result = append(result, RangeInfo{
					Range: overlap,
					Node1: t1.Node,
					Node2: t2.Node,
				})
			}
		}
	}
	return result
}

func findOverlap(range1, range2 []rune) []rune {
	var result []rune
	i, j := 0, 0
	for i < len(range1) && j < len(range2) {
		start := max(range1[i], range2[j])
		end := min(range1[i+1], range2[j+1])

		if start <= end {
			result = append(result, start, end)
		}

		if range1[i+1] < range2[j+1] {
			i += 2
		} else {
			j += 2
		}
	}
	return result
}

func nodeName(node1, node2 *dfa.Node) string {
	return fmt.Sprintf("%d_%d", node1.State, node2.State)
}
