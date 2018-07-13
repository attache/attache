package attache

import "net/http"

// HasRequest is the interface implemented by Context types
// that store an *http.Request
type HasRequest interface {
	Request() *http.Request
	setRequest(*http.Request)
}

// HasResponseWriter is the interface implemented by Context types
// that store an http.ResponseWriter
type HasResponseWriter interface {
	ResponseWriter() http.ResponseWriter
	setResponseWriter(http.ResponseWriter)
}

// HasViews is the interface implemented by Context types
// that use server-rendered views
type HasViews interface {
	CONFIG_Views() ViewConfig
	Views() ViewCache
	setViews(ViewCache)
}

// HasDB is the interaface implemented by Context types
// that use a database connection
type HasDB interface {
	CONFIG_DB() DBConfig
	DB() DB
	setDB(DB)
}

// HasToken is the interface implemented by Context types
// that use managed JWT tokens
type HasToken interface {
	CONFIG_Token() TokenConfig
	Token() Token
	setToken(t Token)
}

// HasSession is the interface implemented by Context types
// that use user sessions
type HasSession interface {
	CONFIG_Session() SessionConfig
	Session() Session
	setSession(s Session)
}
