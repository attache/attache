package attache

type HasViews interface {
	ViewRoot() string
	SetViews(ViewCache)
}

type HasDB interface {
	DBDriver() string
	DBString() string
	SetDB(*DB)
}

type HasFileServer interface {
	FileServer() (dir, path string)
}

type HasMiddleware interface {
	Middleware() Middlewares
}
