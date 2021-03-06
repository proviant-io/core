package http

import (
	"github.com/gorilla/mux"
	"github.com/proviant-io/core/internal/pkg/consumption"
	"github.com/proviant-io/core/internal/pkg/stock"
	"net/http"
	"strconv"
)

func (s *Server) getStock(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
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

	_, customErr := s.productRepo.Get(id, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	models := s.stockRepo.GetAllByProductId(id, accountId)

	var dtos []stock.DTO

	for _, model := range models {
		dtos = append(dtos, stock.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

func (s *Server) addStock(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
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

	dto := stock.DTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	dto.ProductId = id

	if dto.Quantity == 0 {
		response := Response{
			Status: BadRequest,
			Error:  "quantity should not be 0",
		}

		s.jsonResponse(w, response)
		return
	}

	model, customErr := s.relationService.AddStock(dto, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	response := Response{
		Status: ResponseCodeCreated,
		Data:   stock.ModelToDTO(model),
	}

	s.jsonResponse(w, response)
}

func (s *Server) consumeStock(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	userId := s.accountId(r)
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

	dto := stock.ConsumeDTO{}

	err = s.parseJSON(r, &dto)

	if err != nil {
		s.handleBadRequest(w, locale, "parse payload error: %v", err.Error())
		return
	}

	if dto.Quantity == 0 {
		response := Response{
			Status: BadRequest,
			Error:  "quantity should not be 0",
		}

		s.jsonResponse(w, response)
		return
	}

	dto.ProductId = id

	customErr, consumedDTO := s.relationService.ConsumeStock(dto, accountId, userId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	//stock left
	models := s.stockRepo.GetAllByProductId(id, accountId)

	var data struct {
		Stock           []stock.DTO     `json:"stock"`
		ConsumedLogItem consumption.DTO `json:"consumed_log_item"`
	}

	data.Stock = []stock.DTO{}
	data.ConsumedLogItem = consumedDTO

	for _, model := range models {
		data.Stock = append(data.Stock, stock.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   data,
	}

	s.jsonResponse(w, response)
}

func (s *Server) deleteStock(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)
	locale := s.getLocale(r)
	vars := mux.Vars(r)
	productIdString := vars["product_id"]

	if productIdString == "" {
		s.handleBadRequest(w, locale, "product id cannot be empty")
		return
	}
	productId, err := strconv.Atoi(productIdString)

	if err != nil {
		s.handleBadRequest(w, locale, "product id is not a number: %v", err.Error())
		return
	}

	idString := vars["id"]

	if productIdString == "" {
		s.handleBadRequest(w, locale, "id cannot be empty")
		return
	}
	id, err := strconv.Atoi(idString)

	if err != nil {
		s.handleBadRequest(w, locale, "id is not a number: %v", err.Error())
		return
	}

	customErr := s.relationService.DeleteStock(id, accountId)

	if customErr != nil {
		s.handleError(w, locale, *customErr)
		return
	}

	//stock left
	models := s.stockRepo.GetAllByProductId(productId, accountId)

	var dtos []stock.DTO

	for _, model := range models {
		dtos = append(dtos, stock.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}
