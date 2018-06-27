package attache

import (
	"github.com/go-chi/chi"
)

type Router interface {
	chi.Router
}

type router struct {
	chi.Mux
}

type handler interface{}

type node struct {
	prefix string

	list  []handler
	final bool

	skids map[byte]*node
}

func newnode(prefix string, list []handler, final bool) *node {
	return &node{
		prefix: prefix,
		list:   list,
		final:  final,
		skids:  map[byte]*node{},
	}
}

func (n *node) lookup(remaining string) *node {
	if n == nil {
		return nil
	}

	shared := n.sharedPrefix(remaining)

	// matches all of n's prefix?
	if shared == len(n.prefix) {
		// matches all of the remaining path?
		if shared == len(remaining) {
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

func (n *node) insert(path string, list []handler, final bool) {
	shared := n.sharedPrefix(path)

	if shared == len(n.prefix) {
		if shared == len(path) {
			if n.final {
				panic("can't insert to finalized node")
			}

			n.list = append(n.list, list)
			n.final = final
			return
		}

		path = path[shared:]
		if next := n.findChild(path[0]); next != nil {
			next.insert(path, list, final)
			return
		}

		n.addChild(path, list, final)
		return
	}

	n.split(shared)

	if shared == len(path) {
		n.list = append(n.list, list)
		n.final = final
		return
	}

	path = path[shared:]
	n.addChild(path, list, final)
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
	newn := newnode(rest, n.list, n.final)
	newn.skids = n.skids
	n.list = nil
	n.final = false
	n.skids = map[byte]*node{
		rest[0]: newn,
	}
}

func (n *node) findChild(b byte) *node { return n.skids[b] }

func (n *node) addChild(prefix string, list []handler, final bool) {
	label := prefix[0]
	n.skids[label] = newnode(prefix, list, final)
}
