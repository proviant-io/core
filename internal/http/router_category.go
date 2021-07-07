package http

import (
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/brushknight/proviant/internal/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
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

	model, customErr := s.categoryRepo.Get(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   category.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}

func (s *Server) getCategories(w http.ResponseWriter, r *http.Request) {
	models := s.categoryRepo.GetAll()

	dtos := []category.DTO{}

	for _, model := range models {
		dtos = append(dtos, category.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
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

	customErr := s.relationService.DeleteCategory(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
	}

	s.jsonResponse(w, response)
}

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)
	dto := category.DTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Title = utils.ClearString(dto.Title)

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	model := s.categoryRepo.Create(dto)

	response := Response{
		Status: ResponseCodeCreated,
		Data:   category.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
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

	dto := category.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Title = utils.ClearString(dto.Title)

	if dto.Title == "" {
		s.handleBadRequest(w, locale, "title should not be empty")
		return
	}

	model, customErr := s.categoryRepo.Update(id, dto)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   category.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}
