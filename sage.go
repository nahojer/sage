package sage

import (
	"strings"
)

type RouteTrie[T any] struct {
	root *node[T]
}

func NewRouteTrie[T any]() *RouteTrie[T] {
	return &RouteTrie[T]{
		root: &node[T]{},
	}
}

// Add adds value to the trie identified by given HTTP method and route pattern.
// Subsequent calls to Add with the same method and pattern overrides the value.
func (pt *RouteTrie[T]) Add(method, pattern string, value T) {
	segs := pathSegments(strings.TrimRight(pattern, "..."))
	if len(segs) == 0 {
		return
	}

	curr := pt.root
	for _, seg := range segs {
		key := trieKey(method, seg)

		if child, found := curr.children[key]; found {
			curr = child
			continue
		}

		if curr.children == nil {
			curr.children = make(map[string]*node[T])
		}

		toAdd := node[T]{}
		curr.children[key] = &toAdd
		curr = &toAdd
	}

	curr.value = value
	curr.prefix = strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "...")
	curr.valid = true
}

// Lookup searches for the value associated with given HTTP method and URL
// path.
func (pt *RouteTrie[T]) Lookup(method, path string) (value T, found bool) {
	var zero T

	segs := pathSegments(path)
	if len(segs) == 0 {
		return zero, false
	}

	curr := pt.root
	var (
		prefixMatch bool
		prefixValue T
	)
	for _, seg := range segs {
		next, ok := curr.children[trieKey(method, seg)]
		if !ok {
			if prefixMatch {
				break
			}
			return zero, false
		}
		curr = next

		if curr.prefix {
			prefixMatch = true
			prefixValue = curr.value
		}
	}

	if curr.valid {
		return curr.value, true
	}

	if prefixMatch {
		return prefixValue, true
	}

	return zero, false
}

type node[T any] struct {
	children map[string]*node[T]
	valid    bool
	prefix   bool
	value    T
}

func trieKey(method, routeSegment string) string {
	return strings.ToLower(method) + "_" + routeSegment
}

func pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}
