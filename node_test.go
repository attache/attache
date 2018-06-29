package attache

import (
	"fmt"
	"strings"
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
		if c.wantErr != "" {
			if err == nil {
				t.Errorf("%s: wanted error %q, got none", c.path, c.wantErr)
				continue
			}

			if err.Error() != c.wantErr.Error() {
				t.Errorf("%s: wanted error %q, got %q", c.path, c.wantErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("%s: unexpected error %q", c.path, err)
				continue
			}

			if root.lookup(c.path) == nil {
				t.Errorf("%s: did not find newly inserted node", c.path)
			}
		}
	}

	dump(root, "", 0)
}
func dump(root *node, soFar string, deep int) {
	joined := soFar
	if root.prefix != "" {
		joined += root.prefix
	}

	if root.isLeaf() {
		fmt.Println(strings.Repeat(" ", deep), joined, "(leaf)")
	} else {
		methods := []string{}
		for m, _ := range root.methods {
			methods = append(methods, m)
		}
		fmt.Println(strings.Repeat(" ", deep), joined, methods)
	}

	for _, kid := range root.skids {
		dump(kid, joined, deep+1)
	}
}
