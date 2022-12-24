package sage_test

import (
	"testing"

	"github.com/nahojer/sage"
)

var tests = []struct {
	RouteMethod  string
	RoutePattern string
	RouteValue   string

	Method    string
	Path      string
	WantValue string
	Match     bool
}{
	// simple path matching
	{
		"GET", "/one", "one",
		"GET", "/one", "one", true,
	},
	{
		"GET", "/two", "two",
		"GET", "/two", "two", true,
	},
	{
		"GET", "/three", "three",
		"GET", "/three", "three", true,
	},
	// methods
	{
		"get", "/methodcase1", "methodcase1",
		"GET", "/methodcase1", "methodcase1", true,
	},
	{
		"Get", "/methodcase2", "methodcase2",
		"get", "/methodcase2", "methodcase2", true,
	},
	{
		"GET", "/methodcase3", "methodcase3",
		"get", "/methodcase3", "methodcase3", true,
	},
	{
		"GET", "/method1", "method1",
		"POST", "/method1", "", false,
	},
	{
		"DELETE", "/method2", "method2",
		"GET", "/method2", "", false,
	},
	{
		"GET", "/method3", "method3",
		"PUT", "/method3", "", false,
	},
	// nested
	{
		"GET", "/parent/child/one", "nested1",
		"GET", "/parent/child/one", "nested1", true,
	},
	{
		"GET", "/parent/child/two", "nested2",
		"GET", "/parent/child/two", "nested2", true,
	},
	{
		"GET", "/parent/child/three", "nested3",
		"GET", "/parent/child/three", "nested3", true,
	},
	// slashes
	{
		"GET", "slashes/one", "slashes1",
		"GET", "/slashes/one", "slashes1", true,
	},
	{
		"GET", "/slashes/two", "slashes2",
		"GET", "slashes/two", "slashes2", true,
	},
	{
		"GET", "slashes/three/", "slashes3",
		"GET", "/slashes/three", "slashes3", true,
	},
	{
		"GET", "/slashes/four", "slashes4",
		"GET", "slashes/four/", "slashes4", true,
	},
}

func TestRouteTrie(t *testing.T) {
	pt := sage.NewRouteTrie[string]()

	for _, tt := range tests {
		pt.Add(tt.RouteMethod, tt.RoutePattern, tt.RouteValue)
	}

	for _, tt := range tests {
		gotValue, found := pt.Lookup(tt.Method, tt.Path)
		if found != tt.Match || gotValue != tt.WantValue {
			t.Errorf("Lookup(%q, %q) = %q, %t; want %q, %t",
				tt.Method, tt.Path, gotValue, found, tt.WantValue, tt.Match)
		}
	}
}
