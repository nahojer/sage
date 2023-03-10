// Package sage provides support for developing HTTP routers by exporting a
// trie data structure that matches HTTP requests against a list of
// registered routes and returns a route value (typically [http.Handler],
// or a variation of it) for the route that matches the URL and HTTP method.
package sage

import (
	"net/http"
	"strings"
)

// paramKey is the key into nodes that hold parameterized path segments.
const paramKey = "*"

// RoutesTrie is a trie data structure that stores route values of type T.
type RoutesTrie[T any] struct {
	// ParamFunc reports whether given path segment is parameterized and returns
	// the name to give this parameter. The name will be the key into params
	// returned by Lookup.
	//
	// The default ParamFunc consideres a path segment a parameter if it is
	// prefixed with a colon (":"). The returned parameter name is the path
	// segment with all leading colons trimmed.
	ParamFunc func(pathSegment string) (name string, isParam bool)

	root *node[T]
}

// NewRoutesTrie returns a new RoutesTrie.
func NewRoutesTrie[T any]() *RoutesTrie[T] {
	return &RoutesTrie[T]{
		ParamFunc: func(pathSegment string) (name string, isParam bool) {
			if !strings.HasPrefix(pathSegment, ":") {
				return "", false
			}
			return strings.TrimLeft(pathSegment, ":"), true
		},
		root: &node[T]{},
	}
}

// Add inserts a route value to the trie at the location defined by given
// HTTP method and URL path pattern. Subsequent calls to Add with the same
// method and pattern overrides the route value.
//
// Route patterns ending with three dots ("...") are considered prefix routes.
// If there are no matching routes for a HTTP request's URL path and method,
// but a part of the path matches a prefix route, the prefix value will be used.
//
// Path parameters are specified by prefixing a path segment with a colon
// (":"). The parameter name is the value of the path segment with leading
// colons removed. This behaviour can be customized by overriding the ParamFunc
// of the RoutesTrie.
func (rt *RoutesTrie[T]) Add(method, pattern string, value T) {
	method = strings.ToUpper(method)

	segs := pathSegments(strings.TrimRight(pattern, "..."))

	curr := rt.root
	for _, seg := range segs {
		if curr.children == nil {
			curr.children = make(map[string]*node[T])
		}

		key := seg
		if child, found := curr.children[key]; found {
			curr = child
			continue
		}

		var params []string
		if name, isParam := rt.ParamFunc(seg); isParam {
			params = append(params, name)
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

	if curr.values == nil {
		curr.values = make(map[string]T)
	}
	curr.values[method] = value
	curr.prefix = strings.HasSuffix(pattern, "...")
}

// Lookup searches for the route value associated with given HTTP request.
func (rt *RoutesTrie[T]) Lookup(req *http.Request) (value T, params map[string]string, found bool) {
	var zero T

	method := strings.ToUpper(req.Method)

	segs := pathSegments(req.URL.Path)

	curr := rt.root
	var (
		prefixMatch bool
		prefixValue T
	)
	params = make(map[string]string)
	for _, seg := range segs {
		if curr.prefix {
			if value, ok := curr.values[method]; ok {
				prefixMatch = true
				prefixValue = value
			}
		}

		if next, found := curr.children[seg]; found {
			curr = next
			continue
		}

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

	if value, ok := curr.values[method]; ok {
		return value, params, true
	}

	if prefixMatch {
		return prefixValue, params, true
	}

	return zero, nil, false
}

type node[T any] struct {
	children map[string]*node[T]
	params   []string
	prefix   bool
	// All routes values accessed by HTTP method.
	values map[string]T
}

func pathSegments(p string) []string {
	segs := strings.Split(strings.Trim(p, "/"), "/")
	return filter(segs, func(seg string) bool { return seg != "" })
}

func filter[S any](s []S, f func(S) bool) []S {
	var r []S
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}
