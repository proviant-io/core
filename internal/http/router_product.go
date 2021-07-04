package http

import (
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
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

	p, customErr := s.relationService.GetProduct(id)

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
		Data:   p,
	}

	s.JSONResponse(w, response)
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {

	var query *product.Query

	listFilterRaw := r.URL.Query().Get("list")

	listFilter := 0
	var err error

	if listFilterRaw != ""{
		listFilter, err = strconv.Atoi(listFilterRaw)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			s.JSONResponse(w, response)
			return
		}
	}


	categoryFilterRaw := r.URL.Query().Get("category")

	categoryFilter := 0

	if categoryFilterRaw != ""{
		categoryFilter, err = strconv.Atoi(categoryFilterRaw)

		if err != nil {
			response := Response{
				Status: BadRequest,
				Error:  err.Error(),
			}

			s.JSONResponse(w, response)
			return
		}

	}

	if listFilter > 0 || categoryFilter > 0 {
		query = &product.Query{}
		query.List = listFilter
		query.Category = categoryFilter
	}

	dtos := s.relationService.GetAllProducts(query)

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.JSONResponse(w, response)
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request){
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

	customErr := s.relationService.DeleteProduct(id)

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

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request){
	dto := product.CreateDTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	productDto, customErr := s.relationService.CreateProduct(dto)

	if customErr != nil {
		response := Response{
			Status: customErr.Code(),
			Error:  customErr.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	response := Response{
		Status: ResponseCodeCreated,
		Data:   productDto,
	}

	s.JSONResponse(w, response)
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request){
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

	dto := product.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	dto.Id = id

	productDTO, customErr := s.relationService.UpdateProduct(dto)

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
		Data:   productDTO,
	}

	s.JSONResponse(w, response)
}