package attache

// DBConfig provides configuration options for a database connection
type DBConfig struct {
	Driver, DSN string
}

// ViewConfig provides configuration options for initializing
// a ViewCache
type ViewConfig struct {
	Driver, Root string
}

// FileServerConfig provides configuration options for
// starting a file server
type FileServerConfig struct {
	Root, BasePath string
}

// SessionConfig provides configuration options for
// user sessions
type SessionConfig struct {
	Name   string
	Secret []byte
}

// EnvironmentConfig provides configuration options
// for auto-loading environment variables
type EnvironmentConfig struct {
	EnvPath string
}
