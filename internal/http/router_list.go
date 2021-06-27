package http

import (
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		response := Response{
			Status: BadRequest,
			Error:  "id cannot be empty",
		}

		s.JSONResponse(w, response)
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {

		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	model, customErr := s.listRepo.Get(id)

	if customErr != nil {

		response := Response{
			Status: customErr.Code(),
			Error:  customErr.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   list.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}

func (s *Server) getLists(w http.ResponseWriter, r *http.Request) {
	models := s.listRepo.GetAll()

	dtos := []list.DTO{}

	for _, model := range models {
		dtos = append(dtos, list.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   models,
	}

	s.JSONResponse(w, response)
}

func (s *Server) deleteList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		response := Response{
			Status: BadRequest,
			Error:  "id cannot be empty",
		}

		s.JSONResponse(w, response)
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	customErr := s.listRepo.Delete(id)

	if customErr != nil {
		response := Response{
			Status: customErr.Code(),
			Error:  customErr.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
	}

	s.JSONResponse(w, response)
}

func (s *Server) createList(w http.ResponseWriter, r *http.Request) {
	dto := list.DTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	model := s.listRepo.Create(dto)

	response := Response{
		Status: ResponseCodeCreated,
		Data:   list.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}

func (s *Server) updateList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		response := Response{
			Status: BadRequest,
			Error:  "id cannot be empty",
		}

		s.JSONResponse(w, response)
		return
	}
	id, err := strconv.Atoi(idString)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	dto := list.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	model, customErr := s.listRepo.Update(id, dto)

	if customErr != nil {
		response := Response{
			Status: customErr.Code(),
			Error:  customErr.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   list.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}
