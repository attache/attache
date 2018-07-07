package attache

// HasViews is the interface implemented by Context types
// that use server-rendered views
type HasViews interface {
	Context
	CONFIG_Views() ViewConfig
	Views() ViewCache
	SetViews(ViewCache)
}

// HasDB is the interaface implemented by Context types
// that use a database connection
type HasDB interface {
	Context
	CONFIG_DB() DBConfig
	DB() DB
	SetDB(DB)
}

// HasToken is the interface implemented by Context types
// that use managed JWT tokens
type HasToken interface {
	Context
	CONFIG_Token() TokenConfig
	Token() Token
	SetToken(t Token)
}

// HasSession is the interface implemented by Context types
// that use user sessions
type HasSession interface {
	Context
	CONFIG_Session() SessionConfig
	Session() Session
	SetSession(s Session)
}
