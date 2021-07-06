package http

import (
	"net/http"
)

func (s *Server) getMissingTranslations(w http.ResponseWriter, r *http.Request) {

	response := Response{
		Status: ResponseCodeOk,
		Data:   s.l.Missing(),
	}

	s.jsonResponse(w, response)
}
