package http

import (
	"github.com/brushknight/proviant/internal/pkg/list"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getList(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	model, customErr := s.listRepo.Get(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   list.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}

func (s *Server) getLists(w http.ResponseWriter, r *http.Request) {
	models := s.listRepo.GetAll()

	dtos := []list.DTO{}

	for _, model := range models {
		dtos = append(dtos, list.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

func (s *Server) deleteList(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	customErr := s.relationService.DeleteList(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
	}

	s.jsonResponse(w, response)
}

func (s *Server) createList(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)
	dto := list.DTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	model := s.listRepo.Create(dto)

	response := Response{
		Status: ResponseCodeCreated,
		Data:   list.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}

func (s *Server) updateList(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	idString := vars["id"]

	if idString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}
	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	dto := list.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	model, customErr := s.listRepo.Update(id, dto)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   list.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}
