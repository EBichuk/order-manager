package http

import (
	"net/http"

	_ "order-manager/docs"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler *Handler, addr string) *Server {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8081/swagger/doc.json"),
	))

	router.Get("/order/{order_uid}", handler.GetOrder)

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("./pkg/web"))))

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) StartHttpServer() error {
	return s.httpServer.ListenAndServe()
}
