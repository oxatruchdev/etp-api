package http

import (
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api/cmd/web"
)

func (s *Server) registerHomeRoutes() {
	s.Mux.HandleFunc("GET /", s.home)
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	Render(w, r, http.StatusOK, web.Home())
}
