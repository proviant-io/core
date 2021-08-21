package http

import (
	"github.com/proviant-io/core/internal/pkg/shopping"
	"net/http"
)

func (s *Server) getShoppingLists(w http.ResponseWriter, r *http.Request) {
	accountId := s.accountId(r)

	models := s.di.ShoppingList.GetAll(accountId)

	// hack to create first list. Will be removed when multiple shopping lists will be released
	if len(models) == 0 {
		model := s.di.ShoppingList.Create(shopping.ListDTO{
			Title: "Shopping list",
		}, accountId)

		models = append(models, model)
	}

	var dtos []shopping.ListDTO

	for _, model := range models {
		dtos = append(dtos, shopping.ListToDTO(model))
	}

	response := Response{
		Status: ResponseCodeOk,
		Data:   dtos,
	}

	s.jsonResponse(w, response)
}

