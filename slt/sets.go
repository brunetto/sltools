// From
// http://programmers.stackexchange.com/questions/177428/sets-data-structure-in-golang
// http://play.golang.org/p/_FvECoFvhq
// https://groups.google.com/forum/#!topic/golang-nuts/lb4xLHq7wug
// https://github.com/deckarep/golang-set/blob/master/
// Use like:
// set := NewIntSet()
// 	set.Add(1)
// 	set.Add(2)
// 	set.Add(3)
// 	fmt.Println(set.Get(2))
// 	set.Remove(2)
// 	fmt.Println(set.Get(2))
package slt

import (
	"sort"
	"strings"
)

// StringSet is the struct containing the set built as map[string]bool
type StringSet map[string]bool

// NewStringSet creates a new set
// Remember that a map is already a pointer.
func NewStringSet() StringSet {
	return make(StringSet)
}

// NewStringSet creates a new set
func NewStringSetFromSlice(ss []string) StringSet {
	set := NewStringSet()
	for _, item := range ss {
		set.Add(item)
	}
	return set
}

// Add add a new intem to the set if it doesn't exist
func (set StringSet) Add(item string) bool {
	var exists bool
	if _, exists = set[item]; !exists {
		set[item] = true
	}
	return !exists //False if it existed already
}

// Get return the element if it exists
func (set StringSet) Get(item string) bool {
	var exists bool
	_, exists = set[item]
	return exists //true if it existed already
}

// Remove removes an item form the set if it exists
func (set StringSet) Remove(item string) {
	delete(set, item)
}

// Return a sorted slice of string containing the set elements.
func (set StringSet) Sorted() []string {
	temp := make([]string, 0)
	for key, _ := range set {
		temp = append(temp, key)
	}
	sort.Strings(temp)
	return temp
}

// Provides a convenient string representation of the current state of the set.
func (set StringSet) String() string {
	return strings.Join(set.Sorted(), ", ")
}
