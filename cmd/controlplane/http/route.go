package http

import "github.com/resource-aware-jds/resource-aware-jds/pkg/http"

type RouterResult bool

func ProvideHTTPRouter(server http.Server) RouterResult {
	return true
}
