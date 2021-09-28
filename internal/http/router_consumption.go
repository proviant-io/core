package http

import (
	"github.com/gorilla/mux"
	"github.com/proviant-io/core/internal/pkg/consumption"
	"net/http"
	"strconv"
)

func (s *Server) getConsumptionLog(w http.ResponseWriter, r *http.Request) {
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

	models := s.di.ConsumptionLog.GetAllByProductId(id, accountId)

	var dtos []consumption.DTO

	for _, model := range models {
		dtos = append(dtos, consumption.ModelToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

