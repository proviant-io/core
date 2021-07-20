package http

import (
	"github.com/brushknight/proviant/internal/pkg/product"
	"github.com/brushknight/proviant/internal/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getProduct(w http.ResponseWriter, r *http.Request) {
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

	p, customErr := s.relationService.GetProduct(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   p,
	}

	s.jsonResponse(w, response)
}

func (s *Server) getProducts(w http.ResponseWriter, r *http.Request) {
	locale := s.getLocale(r)

	var query *product.Query

	listFilterRaw := r.URL.Query().Get("list")

	listFilter := 0
	var err error

	if listFilterRaw != ""{
		listFilter, err = strconv.Atoi(listFilterRaw)

		if err != nil {
			s.handleBadRequest(w, locale, "list id is not a number: %v", err.Error())
			return
		}
	}


	categoryFilterRaw := r.URL.Query().Get("category")

	categoryFilter := 0

	if categoryFilterRaw != ""{
		categoryFilter, err = strconv.Atoi(categoryFilterRaw)

		if err != nil {
			s.handleBadRequest(w, locale, "category id is not a number: %v", err.Error())
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

	s.jsonResponse(w, response)
}

func (s *Server) deleteProduct(w http.ResponseWriter, r *http.Request){
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

	customErr := s.relationService.DeleteProduct(id)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
	}

	s.jsonResponse(w, response)
}

func (s *Server) createProduct(w http.ResponseWriter, r *http.Request){
	locale := s.getLocale(r)
	dto := product.CreateDTO{}

	err := s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Title = utils.ClearString(dto.Title)
	dto.Image = ""

	productDto, customErr := s.relationService.CreateProduct(dto)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeCreated,
		Data:   productDto,
	}

	s.jsonResponse(w, response)
}

func (s *Server) updateProduct(w http.ResponseWriter, r *http.Request){
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

	dto := product.UpdateDTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.Id = id
	dto.Title = utils.ClearString(dto.Title)

	productDTO, customErr := s.relationService.UpdateProduct(dto)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   productDTO,
	}

	s.jsonResponse(w, response)
}