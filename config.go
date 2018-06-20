package attache

type DBConfig struct {
	Driver, DSN string
}

type ViewConfig struct {
	Root string
}

type FileServerConfig struct {
	Root     string
	BasePath string
}

type TokenConfig struct {
	Secret []byte
	MaxAge int
	Cookie string
}
