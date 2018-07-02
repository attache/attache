package attache

// DBConfig provides configuration options for a database connection
type DBConfig struct {
	Driver, DSN string
}

// ViewConfig provides configuration options for initializing
// a ViewCache
type ViewConfig struct {
	Root string
}

// FileServerConfig provides configuration options for
// starting a file server
type FileServerConfig struct {
	Root     string
	BasePath string
}

// TokenConfig provides configuration options for
// managed JWTs
type TokenConfig struct {
	Secret []byte
	MaxAge int
	Cookie string
}
