package http

type HttpServer interface {
	Listen(port string)
}
