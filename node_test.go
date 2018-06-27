package attache

import (
	"fmt"
	"testing"
)

func TestNodes(t *testing.T) {
	root := newnode("/", nil, false)

	// allow re-insert on non-final node
	try("re-insert on non-final", t, false, func() {
		root.insert("/home", nil, false)
		root.insert("/home", nil, true)
	})

	// forbid re-insert on final node
	try("re-insert on non-final", t, true, func() {
		root.insert("/home", nil, true)
	})

	// remainder
	try("build tree", t, false, func() {
		root.insert("/about", nil, true)
		root.insert("/contact", nil, true)
		root.insert("/contact/form", nil, true)
		root.insert("/register", nil, true)
		root.insert("/item/new", nil, true)
		root.insert("/item/list", nil, true)
		root.insert("/item", nil, true)
		root.insert("/item2/new", nil, true)
		root.insert("/item2/list", nil, true)
		root.insert("/item2", nil, true)
	})

	dump(root, "")
}

func TestNodesFail(t *testing.T) {

}

func try(what string, t *testing.T, wantPanic bool, do func()) {
	defer func() {
		issue := recover()
		if issue == nil && wantPanic {
			t.Errorf("%s: wanted panic, got none", what)
		}

		if issue != nil && !wantPanic {
			t.Errorf("%s: unexpected panic \"%v\"", what, issue)
		}
	}()

	do()
}

func dump(root *node, soFar string) {
	joined := soFar
	if root.prefix != "" {
		if joined != "" {
			joined += "."
		}
		joined += root.prefix
	}
	if root.final {
		fmt.Println(joined)
	} else {
		fmt.Println("[" + joined + "]")
	}
	for _, kid := range root.skids {
		dump(kid, joined)
	}
}
