package http

type Search struct {
	Search string `form:"search"`
	Type   string `form:"type"`
}

func (s *Server) registerSearchRoutes() {
	// s.Mux.Handle("POST /search", s.search)
}

// func (s *Server) search(w http.ResponseWriter, r *http.Request) error {
// 	var search Search
//
// 	if r.Form.Has("search") {
// 		search.Search = r.Form.Get("search")
// 	}
//
// 	if r.Form.Has("type") {
// 		search.Type = r.Form.Get("type")
// 	}
//
// 	slog.Info("Searching", "search", search)
//
// 	return nil
// }
