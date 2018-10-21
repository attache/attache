package attache

import (
	"net/http"
	"reflect"
	"testing"
)

func TestRouter(t *testing.T) {
	cases := []struct {
		method, path string
		wantErr      sentinelError
	}{
		{"GET", "/", ""},
		{"GET", "/a", ""},
		{"GET", "/a/b", ""},
		{"POST", "/a/b", ""},
		{"PUT", "/test", ""},
		{"PUT", "/tent", ""},
		{"GET", "/a/b", errRouteExists},
		{"GET", "/a/b/c/", ""},
		{"", "/web", ""},
		{"", "/web", errRouteExists},
		{"GET", "/we", ""},
		{"GET", "/web/bad", errRoutePastMount},
		{"GET", "/x/y/z", ""},
		{"", "/x/y", errMountOnKnownPath},
	}

	r := &router{
		&node{
			prefix:  "/",
			methods: map[string]stack{},
			kids:    map[byte]*node{},
		},
	}

	for _, c := range cases {
		var err error
		if c.method == "" {
			c.path = canonicalize(c.path, true)
			err = r.mount(c.path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		} else {
			c.path = canonicalize(c.path, false)
			err = r.handle(
				c.method,
				c.path,
				stack{
					reflect.ValueOf(
						func(w http.ResponseWriter, r *http.Request) {},
					),
				},
			)
		}

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

			if r.root.lookup(c.path) == nil {
				t.Errorf("%s: did not find newly inserted node", c.path)
			}
		}
	}

	dump(r.root, "", 0)
}
