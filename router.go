package attache

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"
)

type sentinelError string

func (s sentinelError) String() string { return string(s) }
func (s sentinelError) Error() string  { return string(s) }

const (
	errMountOnKnownPath sentinelError = "illegal mount: path in use by routes"
	errRouteExists      sentinelError = "illegal route: already exists"
	errRoutePastMount   sentinelError = "illegal route: path in use by mount"
)

type router struct {
	root *node
}

func newrouter() router {
	return router{newnode("/", nil, nil)}
}

func (r *router) mount(path string, h http.Handler) {
	path = canonicalize(path, true)
	h = http.StripPrefix(path, h)

	err := r.root.insert("", path, stack{reflect.ValueOf(h.ServeHTTP)}, true)
	if err != nil {
		panic(fmt.Sprintf("mount %s: %s", path, err))
	}
}

func (r *router) handle(method, path string, s stack) {
	path = canonicalize(path, true)
	method = strings.ToUpper(method)

	err := r.root.insert(method, path, s, false)
	if err != nil {
		panic(fmt.Sprintf("route %s: %s", path, err))
	}
}

func (r *router) all(path string, s stack) {
	r.handle("GET", path, s)
	r.handle("PUT", path, s)
	r.handle("POST", path, s)
	r.handle("HEAD", path, s)
	r.handle("TRACE", path, s)
	r.handle("PATCH", path, s)
	r.handle("DELETE", path, s)
	r.handle("OPTIONS", path, s)
}

func canonicalize(p string, trailingSlash bool) string {
	p = path.Join("/", path.Clean(p))
	if trailingSlash {
		return p + "/"
	}
	return p
}

type stack []reflect.Value

type node struct {
	prefix string

	methods map[string]stack
	mounted stack

	skids map[byte]*node
}

func newnode(prefix string, methods map[string]stack, mounted stack) *node {
	n := &node{prefix: prefix}

	if mounted != nil {
		n.mounted = mounted
	} else {
		if methods == nil {
			methods = map[string]stack{}
		}

		n.methods = methods
		n.skids = map[byte]*node{}
	}

	return n
}

func (n *node) lookup(remaining string) *node {
	if n == nil {
		return nil
	}

	shared := n.sharedPrefix(remaining)

	// matches all of n's prefix?
	if shared == len(n.prefix) {
		// matches all of the remaining path?
		if shared == len(remaining) || n.isLeaf() {
			// we've found a match
			return n
		}

		// path remains, try to continue down the trie
		remaining = remaining[shared:]
		return n.findChild(remaining[0]).lookup(remaining)
	}

	// regardless of whether we've matched the whole path remaining,
	// we've fallen in the middle of a node and so we do not have a match
	return nil
}

func (n *node) stackFor(method string) stack {
	if n.isLeaf() {
		return n.mounted
	}

	return n.methods[method]
}

func (n *node) insert(method, path string, s stack, mount bool) error {
	shared := n.sharedPrefix(path)

	if shared == len(n.prefix) {
		if shared == len(path) {
			if mount {
				// it's never valid to mount to a pre-existing node
				return errMountOnKnownPath
			}

			if n.methods[method] != nil {
				return errRouteExists
			}

			n.methods[method] = s
			return nil
		}

		if n.isLeaf() {
			return errRoutePastMount
		}

		path = path[shared:]
		if next := n.findChild(path[0]); next != nil {
			return next.insert(method, path, s, mount)
		}

		n.addChild(method, path, s, mount)
		return nil
	}

	if shared == len(path) {
		if mount {
			// it's never valid to mount to a pre-existing node
			return errMountOnKnownPath
		}

		n.split(shared)
		// otherwise, there's no risk of a redefinition since we just split n
		n.methods[method] = s
		return nil
	}

	n.split(shared)
	path = path[shared:]
	n.addChild(method, path, s, mount)
	return nil
}

func (n *node) sharedPrefix(path string) int {
	max := len(n.prefix)
	if len(path) < len(n.prefix) {
		max = len(path)
	}

	for i := 0; i < max; i++ {
		if n.prefix[i] != path[i] {
			return i
		}
	}
	return max
}

func (n *node) split(at int) {
	var rest string
	n.prefix, rest = n.prefix[:at], n.prefix[at:]
	// create new child
	newn := newnode(rest, n.methods, n.mounted)
	// modify the new parent
	n.methods = map[string]stack{}
	n.mounted = nil
	n.skids = map[byte]*node{
		rest[0]: newn,
	}
}

func (n *node) findChild(b byte) *node { return n.skids[b] }

func (n *node) addChild(method, prefix string, s stack, mount bool) {
	var newn *node
	if mount {
		newn = newnode(prefix, nil, s)
	} else {
		newn = newnode(prefix, map[string]stack{method: s}, nil)
	}
	label := prefix[0]
	n.skids[label] = newn
}

func (n *node) isLeaf() bool { return n.mounted != nil }
