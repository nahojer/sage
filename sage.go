package sage

import (
	"net/http"
	"strings"
)

const paramKey = "*"

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
		if curr.children == nil {
			curr.children = make(map[string]*node[T])
		}

		key := trieKey(method, seg)
		if child, found := curr.children[key]; found {
			curr = child
			continue
		}

		var params []string
		if strings.HasPrefix(seg, ":") {
			params = append(params, strings.TrimLeft(seg, ":"))
		}

		if len(params) > 0 {
			key = paramKey

			if child, found := curr.children[key]; found {
				curr = child
				curr.params = append(curr.params, params...)
				continue
			}
		}

		toAdd := node[T]{params: params}
		curr.children[key] = &toAdd
		curr = &toAdd
	}

	curr.value = value
	curr.prefix = strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "...")
	curr.valid = true
}

// Lookup searches for the value associated with given HTTP method and URL
// path.
func (pt *RouteTrie[T]) Lookup(req *http.Request) (value T, params map[string]string, found bool) {
	var zero T

	segs := pathSegments(req.URL.Path)
	if len(segs) == 0 {
		return zero, nil, false
	}

	curr := pt.root
	var (
		prefixMatch bool
		prefixValue T
	)
	params = make(map[string]string)
	for _, seg := range segs {
		if curr.prefix {
			prefixMatch = true
			prefixValue = curr.value
		}

		key := trieKey(req.Method, seg)

		next, found := curr.children[key]
		if !found {
			if next, found := curr.children[paramKey]; found {
				curr = next
				for _, name := range curr.params {
					params[name] = seg
				}
				continue
			}

			if prefixMatch {
				break
			}

			return zero, nil, false
		}

		curr = next
	}

	if curr.valid {
		return curr.value, params, true
	}

	if prefixMatch {
		return prefixValue, params, true
	}

	return zero, nil, false
}

type node[T any] struct {
	children map[string]*node[T]
	valid    bool
	params   []string
	prefix   bool
	value    T
}

func trieKey(method, routeSegment string) string {
	return strings.ToLower(method) + "_" + routeSegment
}

func pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}
