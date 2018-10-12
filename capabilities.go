package attache

// HasViews is the interface implemented by Context types
// that use server-rendered views
type HasViews interface {
	CONFIG_Views() ViewConfig
	SetViews(ViewCache)
	Views() ViewCache

	ViewData() interface{}
	SetViewData(data interface{})
}

// HasDB is the interaface implemented by Context types
// that use a database connection
type HasDB interface {
	CONFIG_DB() DBConfig
	SetDB(DB)
	DB() DB
}

// HasSession is the interface implemented by Context types
// that use user sessions
type HasSession interface {
	CONFIG_Session() SessionConfig
	SetSession(Session)
	Session() Session
}

// HasEnvironment is the interface implemented by Context types
// that auto-load environment variables
type HasEnvironment interface {
	CONFIG_Environment() EnvironmentConfig
}
