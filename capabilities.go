package attache

// HasViews is the interface implemented by Context types
// that use server-rendered views
type HasViews interface {
	CONFIG_Views() ViewConfig
	Views() ViewCache
	setViews(ViewCache)

	ViewData() interface{}
	SetViewData(data interface{})
}

// HasDB is the interaface implemented by Context types
// that use a database connection
type HasDB interface {
	CONFIG_DB() DBConfig
	DB() DB
	setDB(DB)
}

// HasSession is the interface implemented by Context types
// that use user sessions
type HasSession interface {
	CONFIG_Session() SessionConfig
	Session() Session
	setSession(s Session)
}

// HasEnvironment is the interface implemented by Context types
// that auto-load environment variables
type HasEnvironment interface {
	CONFIG_Environment() EnvironmentConfig
}
