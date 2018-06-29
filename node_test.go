package attache

import (
	"fmt"
	"testing"
)

func TestNodes(t *testing.T) {
	cases := []struct {
		method, path string
		mount        bool
		wantErr      sentinelError
	}{
		{"GET", "/", false, ""},
		{"GET", "/a", false, ""},
		{"GET", "/a/b", false, ""},
		{"POST", "/a/b", false, ""},
		{"PUT", "/test", false, ""},
		{"PUT", "/tent", false, ""},
		{"GET", "/a/b", false, errRouteExists},
		{"GET", "/a/b/c/", false, ""},
		{"", "/web", true, ""},
		{"", "/web", true, errMountOnKnownPath},
		{"GET", "/we", false, ""},
		{"GET", "/web/bad", false, errRoutePastMount},
		{"", "/a/b/c", true, errMountOnKnownPath},
	}

	root := newnode("/", nil, nil)

	for _, c := range cases {
		err := root.insert(c.method, c.path, stack{}, c.mount)
		if err == nil && c.wantErr != "" {
			t.Fatalf("%s: wanted error %q, got none", c.path, c.wantErr)
		}

		if err != nil {
			if c.wantErr == "" {
				t.Fatalf("%s: unexpected error %q", c.path, err)
			}

			if err.Error() != c.wantErr.Error() {
				t.Fatalf("%s: wanted error %q, got %q", c.path, c.wantErr, err)
			}
		}
	}

	dump(root, "")
}
func dump(root *node, soFar string) {
	joined := soFar
	if root.prefix != "" {
		if joined != "" {
			joined += ""
		}
		joined += root.prefix
	}

	if root.isLeaf() {
		fmt.Println(joined, "(leaf)")
	} else {
		methods := []string{}
		for m, _ := range root.methods {
			methods = append(methods, m)
		}
		fmt.Println(joined, methods)
	}

	for _, kid := range root.skids {
		dump(kid, joined)
	}
}
