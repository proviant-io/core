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

func (s *Server) getVersion(w http.ResponseWriter, r *http.Request) {

	type versionResponse struct {
		Version string `json:"version"`
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   versionResponse{
			Version: s.di.Version,
		},
	}

	s.jsonResponse(w, response)
}