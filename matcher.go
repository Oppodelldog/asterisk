package asterisk

import "go/ast"

// New returns a new instance of a Matcher.
func New(conditions []NodeCondition, processMatch func()) *Matcher {
	return &Matcher{conditions: conditions, processMatch: processMatch}
}

// Matcher helps to find ast portions of interest while walking through the tree.
type Matcher struct {
	conditions   []NodeCondition
	idx          int
	processMatch func()
}

// Walk will sequentially called detect node chains by the configured conditions.
// Once all conditions have matched, the ProcessMatch will be called.
func (pm *Matcher) Match(n ast.Node) {
	if pm.conditions[pm.idx](n) {
		pm.idx++
		if pm.idx >= len(pm.conditions) {
			pm.idx = 0
			pm.processMatch()
		}
	} else {
		pm.idx = 0
	}
}

// PatternMatchers holds multiple matchers.
type PatternMatchers []*Matcher

// Walk matches all matchers against the given node.
func (m PatternMatchers) Match(n ast.Node) {
	for _, matcher := range m {
		matcher.Match(n)
	}
}
