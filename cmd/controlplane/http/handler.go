package http

type Handler struct {
	JobHandler JobHandler
}

func ProvideHandler(jobHandler JobHandler) Handler {
	return Handler{
		JobHandler: jobHandler,
	}
}
