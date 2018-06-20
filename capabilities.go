package attache

type HasViews interface {
	CONFIG_Views() ViewConfig
	Views() ViewCache
	SetViews(ViewCache)
}

type HasDB interface {
	CONFIG_DB() DBConfig
	DB() DB
	SetDB(DB)
}

type HasFileServer interface {
	CONFIG_FileServer() FileServerConfig
}

type HasMiddleware interface {
	Middleware() Middlewares
}

type HasToken interface {
	CONFIG_Token() TokenConfig
	Token() Token
	SetToken(t Token)
}
