package http

type Handler struct {
	JobHandler HttpHandler
}

func ProvideHandler(jobHandler HttpHandler) Handler {
	return Handler{
		JobHandler: jobHandler,
	}
}
