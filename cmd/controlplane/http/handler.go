package http

type Handler struct {
	httpHandler HttpHandler
}

func ProvideHandler(jobHandler HttpHandler) Handler {
	return Handler{
		httpHandler: jobHandler,
	}
}
