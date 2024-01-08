func New(cfg Config,
dbRepository handlers.DBRepository,
s3FileRepository handlers.S3FileRepository,
tokenManager handlers.TokenManager,
// merged --->
keyRepository handlers.KeyRepository,
mail *email.Mail,
// <---
logger *slog.Logger) *server {

logger = logger.With(slog.String("from", "http-server"))

srv := &http.Server{
Addr:         fmt.Sprintf(":%d", cfg.Port),
ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
IdleTimeout:  defaultIdleTimeout,
ReadTimeout:  defaultReadTimeout,
WriteTimeout: defaultWriteTimeout,
}

s := &server{
httpServer: srv,
logger:     logger,
}

passwordVerification := password.New()

// TODO: keyRepository, mail перенести в акутальную версию (в виде сервиса?)
handler := handlers.New(dbRepository, s3FileRepository, passwordVerification, tokenManager, keyRepository, mail, logger)

mux := chi.NewRouter()
mux.NotFound(srvErrors.NotFound)
mux.MethodNotAllowed(srvErrors.MethodNotAllowed)
mux.Use(middleware.LogAccess)
mux.Use(middleware.RecoverPanic)

srv.Handler = api.HandlerWithOptions(handler, api.ChiServerOptions{
BaseURL:    baseURL,
BaseRouter: mux,
})

return s
}