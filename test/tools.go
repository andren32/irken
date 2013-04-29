// Package test contains test tools extending the basic go test toolset.
// Based on work by Stefan Nilsson of KTH
package test

import (
	"fmt"
	"sort"
	"testing"
)

// Checks if res and exp are different.
func Diff(res, exp interface{}) (message string, diff bool) {
	switch res := res.(type) {
	case []int:
		if !ArrayEq(res, exp.([]int)) {
			message = fmt.Sprintf("%v; want %v", res, exp)
			diff = true
		}
	default:
		if res != exp {
			message = fmt.Sprintf("%v; want %v", res, exp)
			diff = true
		}
	}
	return
}

// Checks if res and exp are different permutations of characters.
func DiffPerm(res, exp string) (message string, diff bool) {
	r := make([]int, len(res))
	for i, ch := range res {
		r[i] = int(ch)
	}
	e := make([]int, len(exp))
	for i, ch := range exp {
		e[i] = int(ch)
	}
	sort.Ints(r)
	sort.Ints(e)
	if !ArrayEq(r, e) {
		message = fmt.Sprintf("%s; want permutation of %s", res, exp)
		diff = true
	}
	return
}

// Checks if a and b have the same elements in the same order.
func ArrayEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, ai := range a {
		if ai != b[i] {
			return false
		}
	}
	return true
}

func Check(t *testing.T, res, exp interface{}) {
	if mess, diff := Diff(res, exp); diff {
		t.Errorf("%s", mess)
	}

}
