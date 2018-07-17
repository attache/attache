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

type stack []reflect.Value

type node struct {
	prefix string

	guard   stack
	methods map[string]stack
	mount   http.Handler

	kids map[byte]*node
}

func (n *node) isDefined() bool   { return n.hasHandlers() || n.hasMount() }
func (n *node) hasMount() bool    { return n.mount != nil }
func (n *node) hasChildren() bool { return len(n.kids) > 0 }
func (n *node) hasHandlers() bool { return len(n.methods) > 0 }

func (n *node) insert(path string, mustBeLeaf bool) (*node, error) {
	return n.insert2(path, mustBeLeaf, nil)
}

// guards MUST be inserted in order of shortest path to longest
func (n *node) insert2(path string, mustBeLeaf bool, guardStack stack) (*node, error) {
	shared := n.sharedPrefix(path)

	if shared == len(n.prefix) {
		if shared == len(path) {
			if mustBeLeaf && n.hasChildren() {
				return nil, errMountOnKnownPath
			}

			return n, nil
		}

		remaining := path[shared:]
		if next := n.kids[remaining[0]]; next != nil {
			return next.insert2(remaining, mustBeLeaf, guardStack)
		}

		if n.hasMount() {
			return nil, errRoutePastMount
		}

		newn := &node{
			prefix:  remaining,
			guard:   append(guardStack, n.guard...),
			methods: map[string]stack{},
			kids:    map[byte]*node{},
		}

		n.kids[newn.prefix[0]] = newn
		return newn, nil
	}

	if shared == len(path) {
		if mustBeLeaf {
			return nil, errMountOnKnownPath
		}

		n.split(shared)
		return n, nil
	}

	n.split(shared)

	newn := &node{
		prefix:  path[shared:],
		guard:   append(guardStack, n.guard...),
		methods: map[string]stack{},
		kids:    map[byte]*node{},
	}

	n.kids[newn.prefix[0]] = newn
	return newn, nil
}

func (n *node) lookup(remaining string) *node {
	if n == nil {
		return nil
	}

	shared := n.sharedPrefix(remaining)

	// matches all of n's prefix?
	if shared == len(n.prefix) {
		// matches all of the remaining path?
		if shared == len(remaining) || n.hasMount() {
			// we've found a match
			return n
		}

		// path remains, try to continue down the trie
		remaining = remaining[shared:]
		return n.kids[remaining[0]].lookup(remaining)
	}

	// regardless of whether we've matched the whole path remaining,
	// we've stopped in the middle of a node, meaning there's no match
	fmt.Println(n.kids)
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

func (n *node) split(split int) {
	start, end := n.prefix[:split], n.prefix[split:]

	newn := &node{
		prefix:  end,
		guard:   n.guard,
		methods: n.methods,
		mount:   n.mount,
		kids:    n.kids,
	}

	*n = node{
		prefix: start,
		// retain our guard stack
		guard:   newn.guard,
		methods: map[string]stack{},
		kids: map[byte]*node{
			end[0]: newn,
		},
	}
}

type router struct {
	root *node
}

func newrouter() router {
	return router{
		&node{
			prefix:  "/",
			methods: map[string]stack{},
			kids:    map[byte]*node{},
		},
	}
}

func (r *router) mount(path string, h http.Handler) error {
	path = canonicalize(path, true)

	n, err := r.root.insert(path, true)
	if err != nil {
		return err
	}

	if n.hasHandlers() || n.hasMount() {
		return errRouteExists
	}

	n.mount = http.StripPrefix(path, h)
	return nil
}

func (r *router) guard(path string, guards stack) error {
	path = canonicalize(path, false)

	n, err := r.root.insert(path, false)
	if err != nil {
		return err
	}

	n.guard = append(n.guard, guards...)
	return nil
}

func (r *router) handle(method, path string, s stack) error {
	path = canonicalize(path, false)
	method = strings.ToUpper(method)

	n, err := r.root.insert(path, false)
	if err != nil {
		return err
	}

	if n.hasMount() || len(n.methods[method]) > 0 {
		return errRouteExists
	}

	n.methods[method] = s

	return nil
}

func (r *router) all(path string, s stack) error {
	for _, meth := range []string{
		"GET",
		"PUT",
		"POST",
		"HEAD",
		"TRACE",
		"PATCH",
		"DELETE",
		"OPTIONS",
	} {
		if err := r.handle(meth, path, s); err != nil {
			return err
		}
	}

	return nil
}

func canonicalize(p string, trailingSlash bool) string {
	p = path.Join("/", path.Clean(p))
	if trailingSlash {
		return p + "/"
	}
	return p
}

func dump(root *node, soFar string, deep int) {
	const indent = "  "

	joined := soFar
	if root.prefix != "" {
		joined += root.prefix
	}

	if root.hasMount() {
		fmt.Printf("%s%s %s", strings.Repeat(indent, deep), joined, "(mounted)")
	} else {
		fmt.Printf("%s%s", strings.Repeat(indent, deep), joined)
		methods := []string{}
		for m := range root.methods {
			methods = append(methods, m)
		}

		if len(methods) > 0 {
			fmt.Printf(" %v", methods)
		}
	}

	if len(root.guard) > 0 {
		fmt.Printf(" (%d guards)", len(root.guard))
	}

	fmt.Println()

	for sign, kid := range root.kids {
		fmt.Println(strings.Repeat(indent, deep), "child:", string([]byte{sign}))
		dump(kid, joined, deep+1)
	}
}
