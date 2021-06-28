package http

import (
	"github.com/brushknight/proviant/internal/pkg/category"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getCategory(w http.ResponseWriter, r *http.Request) {
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

	model, customErr := s.categoryRepo.Get(id)

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
		Data:   category.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
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

	s.JSONResponse(w, response)
}

func (s *Server) deleteCategory(w http.ResponseWriter, r *http.Request) {
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

	customErr := s.categoryRepo.Delete(id)

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

func (s *Server) createCategory(w http.ResponseWriter, r *http.Request) {
	dto := category.DTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	if dto.Title == "" {
		response := Response{
			Status: BadRequest,
			Error:  "title should not be empty",
		}

		s.JSONResponse(w, response)
		return
	}

	model := s.categoryRepo.Create(dto)

	response := Response{
		Status: ResponseCodeCreated,
		Data:   category.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}

func (s *Server) updateCategory(w http.ResponseWriter, r *http.Request) {
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

	dto := category.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	if dto.Title == "" {
		response := Response{
			Status: BadRequest,
			Error:  "title should not be empty",
		}

		s.JSONResponse(w, response)
		return
	}

	model, customErr := s.categoryRepo.Update(id, dto)

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
		Data:   category.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}


