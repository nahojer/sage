package sage_test

import (
	"net/http/httptest"
	"testing"

	"github.com/nahojer/sage"
)

var tests = []struct {
	RouteMethod  string
	RoutePattern string
	RouteValue   string

	Method     string
	Path       string
	WantValue  string
	WantParams map[string]string
	Match      bool
}{
	// simple path matching
	{
		"GET", "/one", "one",
		"GET", "/one", "one", nil, true,
	},
	{
		"GET", "/two", "two",
		"GET", "/two", "two", nil, true,
	},
	{
		"GET", "/three", "three",
		"GET", "/three", "three", nil, true,
	},
	// methods
	{
		"get", "/methodcase1", "methodcase1",
		"GET", "/methodcase1", "methodcase1", nil, true,
	},
	{
		"Get", "/methodcase2", "methodcase2",
		"get", "/methodcase2", "methodcase2", nil, true,
	},
	{
		"GET", "/methodcase3", "methodcase3",
		"get", "/methodcase3", "methodcase3", nil, true,
	},
	{
		"GET", "/method1", "method1",
		"POST", "/method1", "", nil, false,
	},
	{
		"DELETE", "/method2", "method2",
		"GET", "/method2", "", nil, false,
	},
	{
		"GET", "/method3", "method3",
		"PUT", "/method3", "", nil, false,
	},
	// nested
	{
		"GET", "/parent/child/one", "nested1",
		"GET", "/parent/child/one", "nested1", nil, true,
	},
	{
		"GET", "/parent/child/two", "nested2",
		"GET", "/parent/child/two", "nested2", nil, true,
	},
	{
		"GET", "/parent/child/three", "nested3",
		"GET", "/parent/child/three", "nested3", nil, true,
	},
	// slashes
	{
		"GET", "slashes/one", "slashes1",
		"GET", "/slashes/one", "slashes1", nil, true,
	},
	{
		"GET", "slashes/two/", "slashes2",
		"GET", "/slashes/two", "slashes2", nil, true,
	},
	// prefix
	{
		"GET", "/prefix/", "prefix",
		"GET", "/prefix/anything/else", "prefix", nil, true,
	},
	{
		"GET", "/not-prefix", "not-prefix",
		"GET", "/not-prefix/anything/else", "", nil, false,
	},
	{
		"GET", "/prefixdots...", "prefixdots1",
		"GET", "/prefixdots/anything/else", "prefixdots1", nil, true,
	},
	{
		"POST", "/prefixdots...", "prefixdots2",
		"POST", "/prefixdots", "prefixdots2", nil, true,
	},
	{
		"DELETE", "/prefixdots/...", "prefixdots3",
		"DELETE", "/prefixdots/anything/else", "prefixdots3", nil, true,
	},
	// path params
	{
		"GET", "/path-param/:id", "params1",
		"GET", "/path-param/123", "params1", map[string]string{"id": "123"}, true,
	},
	{
		"GET", "/path-params/:era/:group/:member", "params2",
		"GET", "/path-params/60s/beatles/lennon", "params2", map[string]string{
			"era":    "60s",
			"group":  "beatles",
			"member": "lennon",
		}, true,
	},
	{
		"GET", "/path-params-prefix/:era/:group/:member/", "params3",
		"GET", "/path-params-prefix/60s/beatles/lennon/yoko", "params3", map[string]string{
			"era":    "60s",
			"group":  "beatles",
			"member": "lennon",
		}, true,
	},
	{
		"GET", "/path-params-prefix/:era/:group/award-winners/", "params4",
		"GET", "/path-params-prefix/60s/beatles/award-winners/lennon", "params4", map[string]string{
			"era":   "60s",
			"group": "beatles",
		}, true,
	},
	// misc no matches
	{
		"GET", "/not/enough", "notenough1",
		"GET", "/not/enough/items", "", nil, false,
	},
	{
		"POST", "/not/enough/items", "notenough2",
		"POST", "/not/enough", "", nil, false,
	},
}

func TestRouteTrie(t *testing.T) {
	rt := sage.NewRouteTrie[string]()

	for _, tt := range tests {
		rt.Add(tt.RouteMethod, tt.RoutePattern, tt.RouteValue)
	}

	for _, tt := range tests {
		req := httptest.NewRequest(tt.Method, "http://localhost"+tt.Path, nil)

		gotValue, gotParams, found := rt.Lookup(req)
		if found != tt.Match || gotValue != tt.WantValue || !isSubset(gotParams, tt.WantParams) {
			t.Errorf("Lookup(%q, %q) = %q, %+v, %t; want %q, %+v, %t",
				tt.Method, tt.Path, gotValue, gotParams, found, tt.WantValue, tt.WantParams, tt.Match)
		}
	}
}

// isSubset reports whether sub is a subset of m.
func isSubset[K, V comparable](m, sub map[K]V) bool {
	if len(sub) > len(m) {
		return false
	}
	for k, vsub := range sub {
		if vm, found := m[k]; !found || vm != vsub {
			return false
		}
	}
	return true
}
