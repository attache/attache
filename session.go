package attache

import (
	"github.com/gorilla/sessions"
)

type Session struct {
	*sessions.Session
}
