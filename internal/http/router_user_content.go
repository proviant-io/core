package http

import (
	"github.com/gorilla/mux"
	"github.com/proviant-io/core/internal/errors"
	"github.com/proviant-io/core/internal/i18n"
	"io"
	"log"
	"net/http"
)

func (s *Server) getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imageFileName := vars["fileName"]

	if imageFileName == ""{
		s.handleError(w, s.getLocale(r), *errors.NewInternalServer(i18n.NewMessage("file name cannot be empty")))
		return
	}

	f, err := s.di.ImageSaver.GetImage(imageFileName)

	if err != nil{
		s.handleError(w, s.getLocale(r), *errors.NewInternalServer(i18n.NewMessage("Cannot fetch file, : %s", err.Error())))
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(200)
	_, err = io.Copy(w, f)
	if err != nil {
		log.Println(err)
	}
	return
}
