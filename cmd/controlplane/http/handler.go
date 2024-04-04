package http

type Handler struct {
	httpHandler     HttpHandler
	nodePoolHandler NodeHandler
}

func ProvideHandler(jobHandler HttpHandler, nodePoolHandler NodeHandler) Handler {
	return Handler{
		httpHandler:     jobHandler,
		nodePoolHandler: nodePoolHandler,
	}
}
