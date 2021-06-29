package http

import (
	"github.com/brushknight/proviant/internal/pkg/stock"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (s *Server) getStock(w http.ResponseWriter, r *http.Request) {
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

	models := s.stockRepo.GetAllByProductId(id)

	var dtos []stock.DTO

	for _, model := range models {
		dtos = append(dtos, stock.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.JSONResponse(w, response)
}

func (s *Server) addStock(w http.ResponseWriter, r *http.Request) {
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

	dto := stock.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	dto.ProductId = id

	if dto.Quantity == 0 {
		response := Response{
			Status: BadRequest,
			Error:  "quantity should not be 0",
		}

		s.JSONResponse(w, response)
		return
	}

	model, customErr := s.relationService.AddStock(dto)

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
		Data:   stock.ModelToDTO(model),
	}

	s.JSONResponse(w, response)
}

func (s *Server) consumeStock(w http.ResponseWriter, r *http.Request) {
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

	dto := stock.ConsumeDTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		response := Response{
			Status: BadRequest,
			Error:  err.Error(),
		}

		s.JSONResponse(w, response)
		return
	}

	if dto.Quantity == 0 {
		response := Response{
			Status: BadRequest,
			Error:  "quantity should not be 0",
		}

		s.JSONResponse(w, response)
		return
	}

	dto.ProductId = id

	customErr := s.relationService.ConsumeStock(dto)

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
